"""
Llama 4 Training Script

This script provides functionality for fine-tuning Llama 4 models on
GitOps, Terraform, and Kubernetes issue data.
"""

import os
import json
import logging
import argparse
from pathlib import Path
from typing import Dict, List, Optional, Any

import torch
import mlflow
import numpy as np
from transformers import (
    AutoModelForCausalLM,
    AutoTokenizer,
    Trainer,
    TrainingArguments,
    DataCollatorForLanguageModeling,
    set_seed,
)
from datasets import load_dataset
from peft import LoraConfig, get_peft_model, prepare_model_for_kbit_training

import sys
sys.path.append(str(Path(__file__).parent.parent.parent))

from training.config.llama4_maverick_config import Llama4MaverickConfig, Llama4MaverickTrainingArguments
from training.config.llama4_scout_config import Llama4ScoutConfig, Llama4ScoutTrainingArguments
from training.evaluation.metrics import ModelEvaluationMetrics, calculate_metrics
from training.registry.model_registry import ModelRegistry, ModelRegistryConfig

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)


def parse_args():
    """Parse command line arguments."""
    parser = argparse.ArgumentParser(description="Fine-tune Llama 4 models")
    
    parser.add_argument(
        "--model_type",
        type=str,
        required=True,
        choices=["maverick", "scout"],
        help="Type of Llama 4 model to fine-tune",
    )
    
    parser.add_argument(
        "--dataset_path",
        type=str,
        required=True,
        help="Path to the dataset",
    )
    
    parser.add_argument(
        "--output_dir",
        type=str,
        default=None,
        help="Directory to save the fine-tuned model",
    )
    
    parser.add_argument(
        "--config_path",
        type=str,
        default=None,
        help="Path to the configuration file",
    )
    
    parser.add_argument(
        "--mlflow_tracking_uri",
        type=str,
        default=None,
        help="MLflow tracking URI",
    )
    
    parser.add_argument(
        "--mlflow_experiment_name",
        type=str,
        default=None,
        help="MLflow experiment name",
    )
    
    parser.add_argument(
        "--register_model",
        action="store_true",
        help="Register the model with MLflow Model Registry",
    )
    
    parser.add_argument(
        "--seed",
        type=int,
        default=42,
        help="Random seed",
    )
    
    parser.add_argument(
        "--debug",
        action="store_true",
        help="Enable debug mode",
    )
    
    return parser.parse_args()


def load_config(model_type: str, config_path: Optional[str] = None) -> Dict[str, Any]:
    """Load the configuration for the specified model type."""
    if config_path and os.path.exists(config_path):
        with open(config_path, "r") as f:
            config_dict = json.load(f)
        
        if model_type == "maverick":
            config = Llama4MaverickConfig.from_dict(config_dict)
        else:
            config = Llama4ScoutConfig.from_dict(config_dict)
    else:
        if model_type == "maverick":
            config = Llama4MaverickConfig()
        else:
            config = Llama4ScoutConfig()
    
    return config


def prepare_dataset(dataset_path: str, tokenizer, max_length: int = 2048):
    """Prepare the dataset for fine-tuning."""
    dataset = load_dataset("json", data_files=dataset_path)
    
    def preprocess_function(examples):
        inputs = []
        for repo, title, desc in zip(
            examples["repository"], examples["issue_title"], examples["issue_description"]
        ):
            inputs.append(f"Repository: {repo}\nTitle: {title}\nDescription: {desc}")
        
        model_inputs = tokenizer(
            inputs,
            max_length=max_length,
            truncation=True,
            padding="max_length",
        )
        
        labels = tokenizer(
            examples["solution"],
            max_length=max_length,
            truncation=True,
            padding="max_length",
        ).input_ids
        
        model_inputs["labels"] = labels
        
        return model_inputs
    
    processed_dataset = dataset.map(
        preprocess_function,
        batched=True,
        remove_columns=dataset["train"].column_names,
    )
    
    return processed_dataset


def setup_mlflow(tracking_uri: Optional[str], experiment_name: str):
    """Set up MLflow tracking."""
    if tracking_uri:
        mlflow.set_tracking_uri(tracking_uri)
    
    experiment = mlflow.get_experiment_by_name(experiment_name)
    if experiment is None:
        experiment_id = mlflow.create_experiment(experiment_name)
    else:
        experiment_id = experiment.experiment_id
    
    return experiment_id


def train_model(args):
    """Train the Llama 4 model."""
    set_seed(args.seed)
    
    config = load_config(args.model_type, args.config_path)
    
    if args.output_dir:
        config.output_dir = args.output_dir
    if args.mlflow_tracking_uri:
        config.mlflow_tracking_uri = args.mlflow_tracking_uri
    if args.mlflow_experiment_name:
        config.mlflow_experiment_name = args.mlflow_experiment_name
    
    experiment_id = setup_mlflow(config.mlflow_tracking_uri, config.mlflow_experiment_name)
    
    with mlflow.start_run(experiment_id=experiment_id) as run:
        run_id = run.info.run_id
        logger.info(f"Started MLflow run: {run_id}")
        
        mlflow.log_params(config.to_dict())
        
        tokenizer = AutoTokenizer.from_pretrained(
            config.tokenizer_name,
            revision=config.tokenizer_revision,
            use_fast=True,
        )
        
        dataset = prepare_dataset(args.dataset_path, tokenizer, config.max_seq_length)
        
        model = AutoModelForCausalLM.from_pretrained(
            config.model_name,
            revision=config.model_revision,
            torch_dtype=torch.float16 if config.fp16 else torch.float32,
            device_map="auto",
        )
        
        model = prepare_model_for_kbit_training(model)
        
        if config.use_lora:
            lora_config = LoraConfig(
                r=config.lora_r,
                lora_alpha=config.lora_alpha,
                target_modules=config.lora_target_modules,
                lora_dropout=config.lora_dropout,
                bias="none",
                task_type="CAUSAL_LM",
            )
            model = get_peft_model(model, lora_config)
            
            mlflow.log_params({
                "lora_r": config.lora_r,
                "lora_alpha": config.lora_alpha,
                "lora_dropout": config.lora_dropout,
                "lora_target_modules": str(config.lora_target_modules),
            })
        
        if args.model_type == "maverick":
            training_args = Llama4MaverickTrainingArguments(config).to_transformers_training_arguments()
        else:
            training_args = Llama4ScoutTrainingArguments(config).to_transformers_training_arguments()
        
        training_args["output_dir"] = os.path.join(training_args["output_dir"], run_id)
        
        trainer = Trainer(
            model=model,
            args=TrainingArguments(**training_args),
            train_dataset=dataset["train"],
            eval_dataset=dataset["validation"] if "validation" in dataset else None,
            tokenizer=tokenizer,
            data_collator=DataCollatorForLanguageModeling(tokenizer=tokenizer, mlm=False),
        )
        
        logger.info("Starting training...")
        train_result = trainer.train()
        
        metrics = train_result.metrics
        trainer.log_metrics("train", metrics)
        
        logger.info(f"Saving model to {training_args['output_dir']}")
        trainer.save_model()
        
        if "test" in dataset:
            logger.info("Evaluating model...")
            eval_results = trainer.evaluate(dataset["test"])
            trainer.log_metrics("eval", eval_results)
            
            mlflow.log_metrics(eval_results)
        
        if args.register_model:
            logger.info("Registering model...")
            model_registry = ModelRegistry()
            registry_result = model_registry.register_model(
                model_type=args.model_type,
                run_id=run_id,
                model_path="model",
                metrics=metrics,
                tags={
                    "model_type": args.model_type,
                    "dataset": os.path.basename(args.dataset_path),
                    "run_id": run_id,
                },
            )
            
            logger.info(f"Model registered: {registry_result}")
        
        logger.info(f"Training completed. MLflow run ID: {run_id}")
        
        return {
            "run_id": run_id,
            "metrics": metrics,
            "model_path": training_args["output_dir"],
        }


if __name__ == "__main__":
    args = parse_args()
    
    if args.debug:
        logging.getLogger().setLevel(logging.DEBUG)
    
    train_model(args)

"""
Data preprocessing module for Llama 4 fine-tuning.

This module provides functionality to preprocess data collected from GitHub and Gitee
issue scrapers for fine-tuning Llama 4 models.
"""

import os
import json
import logging
import pandas as pd
import numpy as np
from typing import Dict, List, Optional, Union, Any
from datetime import datetime
import re
import nltk
from nltk.tokenize import word_tokenize
from nltk.corpus import stopwords

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)

try:
    nltk.download('punkt', quiet=True)
    nltk.download('stopwords', quiet=True)
except Exception as e:
    logger.warning(f"Failed to download NLTK resources: {e}")


class DataPreprocessor:
    """
    Data preprocessor class for Llama 4 fine-tuning.
    """

    def __init__(
        self,
        input_dir: str = "data",
        output_dir: str = "processed_data",
        max_input_length: int = 2048,
        max_output_length: int = 2048,
    ):
        """
        Initialize the DataPreprocessor.

        Args:
            input_dir: Directory containing raw data
            output_dir: Directory to store processed data
            max_input_length: Maximum length of input text
            max_output_length: Maximum length of output text
        """
        self.input_dir = input_dir
        self.output_dir = output_dir
        self.max_input_length = max_input_length
        self.max_output_length = max_output_length
        
        os.makedirs(self.output_dir, exist_ok=True)
        
        logger.info("DataPreprocessor initialized")

    def load_data(self, filename: str) -> List[Dict[str, Any]]:
        """
        Load data from a JSON file.

        Args:
            filename: Name of the file to load data from

        Returns:
            List of data examples
        """
        logger.info(f"Loading data from {filename}")
        
        input_path = os.path.join(self.input_dir, filename)
        with open(input_path, "r") as f:
            data = json.load(f)
        
        logger.info(f"Loaded {len(data)} examples from {filename}")
        return data

    def clean_text(self, text: str) -> str:
        """
        Clean text by removing special characters, extra whitespace, etc.

        Args:
            text: Text to clean

        Returns:
            Cleaned text
        """
        if not text:
            return ""
        
        text = re.sub(r'<.*?>', '', text)
        
        text = re.sub(r'http[s]?://\S+', '[URL]', text)
        
        text = re.sub(r'\s+', ' ', text)
        
        text = text.strip()
        
        return text

    def format_for_llama4(
        self, examples: List[Dict[str, Any]], format_type: str = "chat"
    ) -> List[Dict[str, Any]]:
        """
        Format data for Llama 4 fine-tuning.

        Args:
            examples: List of data examples
            format_type: Format type (chat or completion)

        Returns:
            List of formatted examples
        """
        logger.info(f"Formatting {len(examples)} examples for Llama 4 ({format_type})")
        
        formatted_examples = []
        
        for example in examples:
            try:
                input_data = example.get("input", {})
                output_data = example.get("output", {})
                metadata = example.get("metadata", {})
                trajectory = example.get("trajectory", [])
                
                repository = self.clean_text(input_data.get("repository", ""))
                topics = input_data.get("topics", [])
                title = self.clean_text(input_data.get("title", ""))
                description = self.clean_text(input_data.get("description", ""))
                
                solution = self.clean_text(output_data.get("solution", ""))
                
                if not title or not description or not solution:
                    continue
                
                input_text = f"Repository: {repository}\n"
                if topics:
                    input_text += f"Topics: {', '.join(topics)}\n"
                input_text += f"Title: {title}\n"
                input_text += f"Description: {description}"
                
                if len(input_text) > self.max_input_length:
                    input_text = input_text[:self.max_input_length]
                
                if len(solution) > self.max_output_length:
                    solution = solution[:self.max_output_length]
                
                if format_type == "chat":
                    formatted_example = {
                        "messages": [
                            {"role": "system", "content": "You are a helpful assistant that solves software engineering issues."},
                            {"role": "user", "content": input_text},
                            {"role": "assistant", "content": solution},
                        ],
                        "metadata": metadata,
                    }
                else:  # completion
                    formatted_example = {
                        "prompt": input_text,
                        "completion": solution,
                        "metadata": metadata,
                    }
                
                if trajectory:
                    formatted_example["trajectory"] = trajectory
                
                formatted_examples.append(formatted_example)
            except Exception as e:
                logger.error(f"Error formatting example: {e}")
        
        logger.info(f"Formatted {len(formatted_examples)} examples for Llama 4")
        return formatted_examples

    def create_train_val_test_split(
        self,
        data: List[Dict[str, Any]],
        train_ratio: float = 0.8,
        val_ratio: float = 0.1,
        test_ratio: float = 0.1,
        seed: int = 42,
    ) -> Dict[str, List[Dict[str, Any]]]:
        """
        Split data into training, validation, and test sets.

        Args:
            data: List of data examples
            train_ratio: Ratio of training examples
            val_ratio: Ratio of validation examples
            test_ratio: Ratio of test examples
            seed: Random seed

        Returns:
            Dictionary containing training, validation, and test sets
        """
        logger.info(f"Splitting {len(data)} examples into train/val/test sets")
        
        np.random.seed(seed)
        
        indices = np.random.permutation(len(data))
        
        train_idx = int(len(data) * train_ratio)
        val_idx = int(len(data) * (train_ratio + val_ratio))
        
        train_data = [data[i] for i in indices[:train_idx]]
        val_data = [data[i] for i in indices[train_idx:val_idx]]
        test_data = [data[i] for i in indices[val_idx:]]
        
        logger.info(f"Split data into {len(train_data)} training, {len(val_data)} validation, and {len(test_data)} test examples")
        
        return {
            "train": train_data,
            "val": val_data,
            "test": test_data,
        }

    def save_processed_data(
        self, data: Dict[str, List[Dict[str, Any]]], prefix: str = "llama4"
    ) -> Dict[str, str]:
        """
        Save processed data to files.

        Args:
            data: Dictionary containing training, validation, and test sets
            prefix: Prefix for output filenames

        Returns:
            Dictionary containing paths to saved files
        """
        logger.info(f"Saving processed data with prefix {prefix}")
        
        os.makedirs(self.output_dir, exist_ok=True)
        
        output_paths = {}
        for split, examples in data.items():
            output_filename = f"{prefix}_{split}.json"
            output_path = os.path.join(self.output_dir, output_filename)
            
            with open(output_path, "w") as f:
                json.dump(examples, f, indent=2)
            
            output_paths[split] = output_path
            logger.info(f"Saved {len(examples)} {split} examples to {output_path}")
        
        return output_paths

    def generate_data_statistics(
        self, data: Dict[str, List[Dict[str, Any]]]
    ) -> Dict[str, Dict[str, Any]]:
        """
        Generate statistics for processed data.

        Args:
            data: Dictionary containing training, validation, and test sets

        Returns:
            Dictionary containing statistics for each split
        """
        logger.info("Generating data statistics")
        
        statistics = {}
        
        for split, examples in data.items():
            split_stats = {
                "count": len(examples),
                "input_length": {
                    "mean": 0,
                    "median": 0,
                    "min": 0,
                    "max": 0,
                },
                "output_length": {
                    "mean": 0,
                    "median": 0,
                    "min": 0,
                    "max": 0,
                },
                "sources": {},
                "repositories": {},
            }
            
            input_lengths = []
            output_lengths = []
            sources = {}
            repositories = {}
            
            for example in examples:
                if "messages" in example:  # chat format
                    input_text = example["messages"][1]["content"]
                    output_text = example["messages"][2]["content"]
                else:  # completion format
                    input_text = example["prompt"]
                    output_text = example["completion"]
                
                input_lengths.append(len(input_text))
                output_lengths.append(len(output_text))
                
                metadata = example.get("metadata", {})
                source = metadata.get("source", "unknown")
                repository = metadata.get("repository", "unknown")
                
                sources[source] = sources.get(source, 0) + 1
                repositories[repository] = repositories.get(repository, 0) + 1
            
            if input_lengths:
                split_stats["input_length"]["mean"] = np.mean(input_lengths)
                split_stats["input_length"]["median"] = np.median(input_lengths)
                split_stats["input_length"]["min"] = np.min(input_lengths)
                split_stats["input_length"]["max"] = np.max(input_lengths)
            
            if output_lengths:
                split_stats["output_length"]["mean"] = np.mean(output_lengths)
                split_stats["output_length"]["median"] = np.median(output_lengths)
                split_stats["output_length"]["min"] = np.min(output_lengths)
                split_stats["output_length"]["max"] = np.max(output_lengths)
            
            split_stats["sources"] = sources
            split_stats["repositories"] = repositories
            
            statistics[split] = split_stats
        
        logger.info("Generated data statistics")
        return statistics

    def save_statistics(
        self, statistics: Dict[str, Dict[str, Any]], filename: str = "statistics.json"
    ) -> str:
        """
        Save statistics to a file.

        Args:
            statistics: Dictionary containing statistics
            filename: Name of the file to save statistics to

        Returns:
            Path to the saved file
        """
        logger.info(f"Saving statistics to {filename}")
        
        os.makedirs(self.output_dir, exist_ok=True)
        
        output_path = os.path.join(self.output_dir, filename)
        with open(output_path, "w") as f:
            json.dump(statistics, f, indent=2)
        
        logger.info(f"Statistics saved to {output_path}")
        return output_path

    def preprocess_data(
        self,
        input_filename: str,
        output_prefix: str = "llama4",
        format_type: str = "chat",
        train_ratio: float = 0.8,
        val_ratio: float = 0.1,
        test_ratio: float = 0.1,
        seed: int = 42,
    ) -> Dict[str, str]:
        """
        Preprocess data for Llama 4 fine-tuning.

        Args:
            input_filename: Name of the input file
            output_prefix: Prefix for output filenames
            format_type: Format type (chat or completion)
            train_ratio: Ratio of training examples
            val_ratio: Ratio of validation examples
            test_ratio: Ratio of test examples
            seed: Random seed

        Returns:
            Dictionary containing paths to saved files
        """
        logger.info(f"Preprocessing data from {input_filename}")
        
        data = self.load_data(input_filename)
        
        formatted_data = self.format_for_llama4(data, format_type=format_type)
        
        split_data = self.create_train_val_test_split(
            formatted_data,
            train_ratio=train_ratio,
            val_ratio=val_ratio,
            test_ratio=test_ratio,
            seed=seed,
        )
        
        output_paths = self.save_processed_data(split_data, prefix=output_prefix)
        
        statistics = self.generate_data_statistics(split_data)
        stats_path = self.save_statistics(
            statistics, filename=f"{output_prefix}_statistics.json"
        )
        output_paths["statistics"] = stats_path
        
        logger.info(f"Preprocessing completed: {output_paths}")
        return output_paths


def main():
    """
    Main function for testing the DataPreprocessor.
    """
    preprocessor = DataPreprocessor(
        input_dir="data",
        output_dir="processed_data",
        max_input_length=2048,
        max_output_length=2048,
    )
    
    output_paths = preprocessor.preprocess_data(
        input_filename="github_gitee_issues.json",
        output_prefix="llama4",
        format_type="chat",
        train_ratio=0.8,
        val_ratio=0.1,
        test_ratio=0.1,
        seed=42,
    )
    
    logger.info(f"Preprocessing completed: {output_paths}")


if __name__ == "__main__":
    main()

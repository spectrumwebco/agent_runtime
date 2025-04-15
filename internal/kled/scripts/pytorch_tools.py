"""
PyTorch tools for the Kled.io framework.
"""

import os
import json
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import DataLoader, Dataset
import torchvision.transforms as transforms
import numpy as np


class SimpleDataset(Dataset):
    """Simple dataset for PyTorch."""

    def __init__(self, data, transform=None):
        """Initialize dataset."""
        self.data = data
        self.transform = transform

    def __len__(self):
        """Return length of dataset."""
        return len(self.data)

    def __getitem__(self, idx):
        """Get item from dataset."""
        item = self.data[idx]
        if self.transform:
            item = self.transform(item)
        return item


class SimpleModel(nn.Module):
    """Simple model for PyTorch."""

    def __init__(self, input_size, hidden_size, output_size):
        """Initialize model."""
        super(SimpleModel, self).__init__()
        self.layer1 = nn.Linear(input_size, hidden_size)
        self.relu = nn.ReLU()
        self.layer2 = nn.Linear(hidden_size, output_size)

    def forward(self, x):
        """Forward pass."""
        x = self.layer1(x)
        x = self.relu(x)
        x = self.layer2(x)
        return x


def run_inference(model_name, input_data):
    """
    Run inference with the specified model and input.

    Args:
        model_name: Name of the model
        input_data: Input data for inference

    Returns:
        Inference result
    """
    try:
        model_path = os.path.join(os.environ.get("MODELS_DIR", "./models"), f"{model_name}.pt")
        if not os.path.exists(model_path):
            return {"error": f"Model {model_name} not found"}

        model = torch.load(model_path)
        model.eval()

        input_tensor = torch.tensor(input_data["data"], dtype=torch.float32)

        with torch.no_grad():
            output = model(input_tensor)

        result = output.numpy().tolist()

        return {"result": result}
    except Exception as e:
        return {"error": str(e)}


def train_model(model_name, config):
    """
    Train a model with the specified configuration.

    Args:
        model_name: Name of the model
        config: Training configuration

    Returns:
        Training result
    """
    try:
        input_size = config.get("input_size", 10)
        hidden_size = config.get("hidden_size", 20)
        output_size = config.get("output_size", 2)
        model = SimpleModel(input_size, hidden_size, output_size)

        optimizer = optim.Adam(model.parameters(), lr=config.get("learning_rate", 0.001))

        loss_fn = nn.MSELoss()

        data = config.get("data", [])
        dataset = SimpleDataset(data)
        dataloader = DataLoader(dataset, batch_size=config.get("batch_size", 32), shuffle=True)

        epochs = config.get("epochs", 10)
        losses = []
        for epoch in range(epochs):
            epoch_loss = 0
            for batch in dataloader:
                optimizer.zero_grad()
                output = model(batch["input"])
                loss = loss_fn(output, batch["target"])
                loss.backward()
                optimizer.step()
                epoch_loss += loss.item()
            losses.append(epoch_loss / len(dataloader))

        model_path = os.path.join(os.environ.get("MODELS_DIR", "./models"), f"{model_name}.pt")
        torch.save(model, model_path)

        return {"result": "success", "losses": losses}
    except Exception as e:
        return {"error": str(e)}


def save_model(model_name, path):
    """
    Save a model to the specified path.

    Args:
        model_name: Name of the model
        path: Path to save the model


    Returns:
        Save result
    """
    try:
        model_path = os.path.join(os.environ.get("MODELS_DIR", "./models"), f"{model_name}.pt")
        if not os.path.exists(model_path):
            return {"error": f"Model {model_name} not found"}

        model = torch.load(model_path)

        torch.save(model, path)

        return {"result": "success"}
    except Exception as e:
        return {"error": str(e)}


def load_model(model_name, path):
    """
    Load a model from the specified path.

    Args:
        model_name: Name of the model
        path: Path to load the model from

    Returns:
        Load result
    """
    try:
        model = torch.load(path)

        model_path = os.path.join(os.environ.get("MODELS_DIR", "./models"), f"{model_name}.pt")
        torch.save(model, model_path)

        return {"result": "success"}
    except Exception as e:
        return {"error": str(e)}


def evaluate_model(model_name, data):
    """
    Evaluate a model with the specified data.

    Args:
        model_name: Name of the model
        data: Evaluation data

    Returns:
        Evaluation result
    """
    try:
        model_path = os.path.join(os.environ.get("MODELS_DIR", "./models"), f"{model_name}.pt")
        if not os.path.exists(model_path):
            return {"error": f"Model {model_name} not found"}

        model = torch.load(model_path)
        model.eval()

        dataset = SimpleDataset(data)
        dataloader = DataLoader(dataset, batch_size=32, shuffle=False)

        loss_fn = nn.MSELoss()

        total_loss = 0
        with torch.no_grad():
            for batch in dataloader:
                output = model(batch["input"])
                loss = loss_fn(output, batch["target"])
                total_loss += loss.item()

        return {"result": "success", "loss": total_loss / len(dataloader)}
    except Exception as e:
        return {"error": str(e)}

"""
ML Infrastructure API Client

This module provides a client for interacting with the ML infrastructure for Llama 4 fine-tuning.
"""

import os
import json
import time
import logging
import requests
from typing import Dict, List, Any, Optional, Union

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

class MLInfrastructureClient:
    """Client for interacting with the ML infrastructure for Llama 4 fine-tuning."""

    def __init__(
        self,
        base_url: Optional[str] = None,
        username: Optional[str] = None,
        password: Optional[str] = None,
        token: Optional[str] = None,
        token_expiry: Optional[int] = None,
    ):
        """Initialize the ML infrastructure client."""
        self.base_url = base_url or os.environ.get("ML_API_BASE_URL")
        self.username = username or os.environ.get("ML_API_USERNAME")
        self.password = password or os.environ.get("ML_API_PASSWORD")
        self.token = token
        self.token_expiry = token_expiry or 0
        self.max_retries = 3
        self.retry_backoff = 2  # seconds

    def _get_auth_token(self) -> str:
        """Get or refresh the authentication token."""
        if self.token and time.time() < self.token_expiry:
            return self.token

        try:
            auth_response = requests.post(
                f"{self.base_url}/token",
                data={
                    "username": self.username,
                    "password": self.password,
                    "grant_type": "password",
                },
                timeout=10,
            )

            auth_response.raise_for_status()
            token_data = auth_response.json()
            self.token = token_data["access_token"]
            self.token_expiry = time.time() + token_data["expires_in"]
            return self.token

        except requests.exceptions.RequestException as e:
            logger.error(f"Failed to get authentication token: {e}")
            raise Exception(f"Authentication failed: {e}")

    def _get_headers(self) -> Dict[str, str]:
        """Get headers with authentication token."""
        token = self._get_auth_token()
        return {
            "Authorization": f"Bearer {token}",
            "Content-Type": "application/json",
        }

    def _make_request(
        self,
        method: str,
        endpoint: str,
        data: Optional[Dict] = None,
        params: Optional[Dict] = None,
        files: Optional[Dict] = None,
        timeout: int = 30,
    ) -> Any:
        """Make an API request with error handling and retries."""
        url = f"{self.base_url}{endpoint}"
        headers = self._get_headers() if endpoint != "/status" else {"Content-Type": "application/json"}

        retry_count = 0
        last_exception = None

        while retry_count < self.max_retries:
            try:
                if files:
                    headers.pop("Content-Type", None)

                response = requests.request(
                    method=method,
                    url=url,
                    headers=headers,
                    json=data,
                    params=params,
                    files=files,
                    timeout=timeout,
                )

                response.raise_for_status()
                if response.content:
                    return response.json()
                return {}

            except requests.exceptions.HTTPError as e:
                last_exception = e
                status_code = e.response.status_code

                if status_code == 401:
                    self.token = None
                    retry_count += 1
                    continue
                elif status_code == 429:
                    wait_time = self.retry_backoff * (2 ** retry_count)
                    time.sleep(wait_time)
                    retry_count += 1
                    continue
                else:
                    logger.error(f"HTTP error: {e}")
                    raise

            except (requests.exceptions.ConnectionError, requests.exceptions.Timeout) as e:
                last_exception = e
                wait_time = self.retry_backoff * (2 ** retry_count)
                time.sleep(wait_time)
                retry_count += 1
                continue

            except requests.exceptions.RequestException as e:
                logger.error(f"Request error: {e}")
                raise

        logger.error(f"Failed after {self.max_retries} retries")
        raise last_exception or Exception(f"Failed after {self.max_retries} retries")

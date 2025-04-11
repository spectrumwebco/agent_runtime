import unittest
import os
import sys
import json
from unittest.mock import patch, MagicMock

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../../../src')))

from integrations.gitee.scraper import GiteeScraper
from integrations.issue_collector.collector import IssueCollector


class TestGiteeScraper(unittest.TestCase):
    """Integration tests for the Gitee scraper component."""

    def setUp(self):
        """Set up test fixtures."""
        self.scraper = GiteeScraper(
            token=os.environ.get('GITEE_TOKEN', 'dummy_token'),
            max_issues=5
        )
        self.collector = IssueCollector()
        
        self.test_data_dir = os.path.join(os.path.dirname(__file__), '../../data/test')
        os.makedirs(self.test_data_dir, exist_ok=True)

    @patch('integrations.gitee.scraper.GiteeScraper._fetch_issues')
    def test_scraper_integration_with_collector(self, mock_fetch_issues):
        """Test that the Gitee scraper integrates correctly with the issue collector."""
        mock_issues = [
            {
                "number": 1,
                "title": "修复 Kubernetes 部署问题",
                "body": "部署失败，错误：ResourceQuotaExceeded",
                "state": "closed",
                "labels": [{"name": "bug"}, {"name": "kubernetes"}],
                "created_at": "2023-01-01T00:00:00Z",
                "closed_at": "2023-01-02T00:00:00Z",
                "user": {"login": "user1"},
                "html_url": "https://gitee.com/org/repo/issues/1"
            },
            {
                "number": 2,
                "title": "更新 Terraform 模块以支持最新的 AWS 提供程序",
                "body": "需要更新模块以支持最新的 AWS 提供程序版本",
                "state": "closed",
                "labels": [{"name": "enhancement"}, {"name": "terraform"}],
                "created_at": "2023-01-03T00:00:00Z",
                "closed_at": "2023-01-04T00:00:00Z",
                "user": {"login": "user2"},
                "html_url": "https://gitee.com/org/repo/issues/2"
            }
        ]
        
        mock_comments = {
            1: [
                {"user": {"login": "user3"}, "body": "我通过增加资源配额解决了这个问题。"},
                {"user": {"login": "user1"}, "body": "谢谢，这样解决了！"}
            ],
            2: [
                {"user": {"login": "user4"}, "body": "已将模块更新为使用 AWS 提供程序 4.0。"},
                {"user": {"login": "user2"}, "body": "看起来不错，现在合并。"}
            ]
        }
        
        mock_fetch_issues.return_value = mock_issues
        self.scraper._fetch_comments = MagicMock(side_effect=lambda issue_number: mock_comments.get(issue_number, []))
        
        issues = self.scraper.scrape_issues("kubernetes/kubernetes")
        processed_data = self.collector.process_issues(issues, source="gitee", repo="kubernetes/kubernetes")
        
        self.assertEqual(len(processed_data), 2)
        self.assertEqual(processed_data[0]["issue_number"], 1)
        self.assertEqual(processed_data[0]["title"], "修复 Kubernetes 部署问题")
        self.assertEqual(processed_data[0]["source"], "gitee")
        self.assertEqual(processed_data[0]["repository"], "kubernetes/kubernetes")
        self.assertIn("kubernetes", processed_data[0]["labels"])
        
        with open(os.path.join(self.test_data_dir, 'gitee_processed_data.json'), 'w') as f:
            json.dump(processed_data, f, indent=2)

    @patch('integrations.gitee.scraper.GiteeScraper._fetch_issues')
    def test_filter_issues_by_label(self, mock_fetch_issues):
        """Test that issues can be filtered by label."""
        mock_issues = [
            {
                "number": 1,
                "title": "修复 Kubernetes 部署问题",
                "body": "部署失败，错误：ResourceQuotaExceeded",
                "state": "closed",
                "labels": [{"name": "bug"}, {"name": "kubernetes"}],
                "created_at": "2023-01-01T00:00:00Z",
                "closed_at": "2023-01-02T00:00:00Z",
                "user": {"login": "user1"},
                "html_url": "https://gitee.com/org/repo/issues/1"
            },
            {
                "number": 2,
                "title": "更新 Terraform 模块以支持最新的 AWS 提供程序",
                "body": "需要更新模块以支持最新的 AWS 提供程序版本",
                "state": "closed",
                "labels": [{"name": "enhancement"}, {"name": "terraform"}],
                "created_at": "2023-01-03T00:00:00Z",
                "closed_at": "2023-01-04T00:00:00Z",
                "user": {"login": "user2"},
                "html_url": "https://gitee.com/org/repo/issues/2"
            }
        ]
        
        mock_fetch_issues.return_value = mock_issues
        self.scraper._fetch_comments = MagicMock(return_value=[])
        
        issues = self.scraper.scrape_issues("kubernetes/kubernetes", labels=["kubernetes"])
        
        self.assertEqual(len(issues), 1)
        self.assertEqual(issues[0]["number"], 1)
        self.assertEqual(issues[0]["title"], "修复 Kubernetes 部署问题")

    def test_data_format_validation(self):
        """Test that the data format is valid for training."""
        data_file = os.path.join(self.test_data_dir, 'gitee_processed_data.json')
        if not os.path.exists(data_file):
            self.skipTest("Processed data file not found. Run test_scraper_integration_with_collector first.")
        
        with open(data_file, 'r') as f:
            processed_data = json.load(f)
        
        for item in processed_data:
            self.assertIn("issue_number", item)
            self.assertIn("title", item)
            self.assertIn("body", item)
            self.assertIn("state", item)
            self.assertIn("labels", item)
            self.assertIn("created_at", item)
            self.assertIn("closed_at", item)
            self.assertIn("user", item)
            self.assertIn("url", item)
            self.assertIn("comments", item)
            self.assertIn("source", item)
            self.assertIn("repository", item)


if __name__ == '__main__':
    unittest.main()

"""
Static analysis tool for the Veigar cybersecurity agent.

This module provides static code analysis capabilities for the Veigar agent,
integrating tools from various cybersecurity repositories.
"""

import os
import logging
import random
from typing import Any, Dict, List

logger = logging.getLogger(__name__)


class StaticAnalysisTool:
    """Static code analysis tool for security vulnerabilities."""

    def __init__(self):
        """Initialize the static analysis tool."""
        self.tools = self._initialize_tools()
        logger.info("Initialized static analysis tool with %d analyzers", len(self.tools))

    def _initialize_tools(self) -> List[Dict[str, Any]]:
        """Initialize the static analysis tools."""
        return [
            {
                "name": "semgrep",
                "description": "Lightweight static analysis for many languages",
                "languages": ["python", "javascript", "go", "java", "c", "cpp"],
                "source": "awesome-static-analysis",
                "enabled": True
            },
            {
                "name": "bandit",
                "description": "Security oriented static analyzer for Python code",
                "languages": ["python"],
                "source": "awesome-static-analysis",
                "enabled": True
            },
            {
                "name": "gosec",
                "description": "Go security checker",
                "languages": ["go"],
                "source": "awesome-static-analysis",
                "enabled": True
            },
            {
                "name": "eslint-security",
                "description": "ESLint plugin for security linting in JavaScript",
                "languages": ["javascript"],
                "source": "awesome-static-analysis",
                "enabled": True
            },
            {
                "name": "cppcheck",
                "description": "Static analysis tool for C/C++ code",
                "languages": ["c", "cpp"],
                "source": "awesome-static-analysis",
                "enabled": True
            },
            {
                "name": "sonarqube",
                "description": "Continuous code quality and security platform",
                "languages": ["python", "javascript", "go", "java", "c", "cpp"],
                "source": "awesome-static-analysis",
                "enabled": True
            },
            {
                "name": "flawfinder",
                "description": "Examines C/C++ source code for security flaws",
                "languages": ["c", "cpp"],
                "source": "awesome-static-analysis",
                "enabled": True
            },
            {
                "name": "brakeman",
                "description": "Static analysis security vulnerability scanner for Ruby on Rails",
                "languages": ["ruby"],
                "source": "awesome-static-analysis",
                "enabled": True
            }
        ]

    def analyze(self, repository: str, branch: str, files: List[str]) -> Dict[str, Any]:
        """
        Perform static analysis on the specified files.

        Args:
            repository: Repository name
            branch: Branch name
            files: List of files to analyze

        Returns:
            Dict: Static analysis results
        """
        logger.info("Performing static analysis on %d files in %s:%s", len(files), repository, branch)

        files_by_language = self._group_files_by_language(files)

        results = {
            "status": "success",
            "repository": repository,
            "branch": branch,
            "findings": []
        }

        for language, language_files in files_by_language.items():
            language_tools = self._get_tools_for_language(language)

            if not language_tools:
                logger.warning("No static analysis tools available for %s", language)
                continue

            for tool in language_tools:
                try:
                    tool_results = self._run_tool(tool, language_files, repository, branch)
                    results["findings"].extend(tool_results)
                except Exception as e:
                    logger.error("Error running %s: %s", tool['name'], e)
                    results["findings"].append({
                        "tool": tool["name"],
                        "status": "error",
                        "error": str(e)
                    })

        results["findings"] = self._deduplicate_findings(results["findings"])

        results["summary"] = {
            "total_findings": len(results["findings"]),
            "critical": len([f for f in results["findings"] if f.get("severity") == "critical"]),
            "high": len([f for f in results["findings"] if f.get("severity") == "high"]),
            "medium": len([f for f in results["findings"] if f.get("severity") == "medium"]),
            "low": len([f for f in results["findings"] if f.get("severity") == "low"]),
            "info": len([f for f in results["findings"] if f.get("severity") == "info"])
        }

        logger.info("Static analysis complete with %d findings", 
                   results['summary']['total_findings'])

        return results

    def _group_files_by_language(self, files: List[str]) -> Dict[str, List[str]]:
        """Group files by language based on file extension."""
        extensions_map = {
            ".py": "python",
            ".js": "javascript",
            ".ts": "typescript",
            ".go": "go",
            ".java": "java",
            ".c": "c",
            ".cpp": "cpp",
            ".h": "c",
            ".hpp": "cpp",
            ".rb": "ruby",
            ".php": "php",
            ".cs": "csharp",
            ".swift": "swift",
            ".kt": "kotlin",
            ".rs": "rust"
        }

        files_by_language = {}

        for file in files:
            ext = os.path.splitext(file)[1].lower()
            language = extensions_map.get(ext)

            if language:
                if language not in files_by_language:
                    files_by_language[language] = []
                files_by_language[language].append(file)

        return files_by_language

    def _get_tools_for_language(self, language: str) -> List[Dict[str, Any]]:
        """Get tools that support the specified language."""
        return [
            tool for tool in self.tools
            if tool["enabled"] and language in tool["languages"]
        ]

    def _run_tool(
        self,
        tool: Dict[str, Any],
        files: List[str],
        repository: str,
        branch: str
    ) -> List[Dict[str, Any]]:
        """
        Run a static analysis tool on the specified files.

        In a real implementation, this would execute the actual tool.
        For now, we'll simulate tool execution with realistic findings.
        """
        tool_name = tool["name"]
        findings = []

        if tool_name == "semgrep":
            findings = self._simulate_semgrep_findings(files)
        elif tool_name == "bandit":
            findings = self._simulate_bandit_findings(files)
        elif tool_name == "gosec":
            findings = self._simulate_gosec_findings(files)
        elif tool_name == "eslint-security":
            findings = self._simulate_eslint_findings(files)
        elif tool_name == "cppcheck":
            findings = self._simulate_cppcheck_findings(files)
        elif tool_name == "sonarqube":
            findings = self._simulate_sonarqube_findings(files)
        elif tool_name == "flawfinder":
            findings = self._simulate_flawfinder_findings(files)
        elif tool_name == "brakeman":
            findings = self._simulate_brakeman_findings(files)

        for finding in findings:
            finding["tool"] = tool_name
            finding["source"] = tool["source"]

        return findings

    def _simulate_semgrep_findings(self, files: List[str]) -> List[Dict[str, Any]]:
        """Simulate semgrep findings."""
        findings = []

        issues = [
            {
                "title": "SQL Injection",
                "description": "Potential SQL injection vulnerability detected",
                "severity": "high",
                "cwe": "CWE-89",
                "remediation": "Use parameterized queries or prepared statements"
            },
            {
                "title": "Cross-Site Scripting (XSS)",
                "description": "Potential XSS vulnerability detected",
                "severity": "high",
                "cwe": "CWE-79",
                "remediation": "Use context-specific output encoding"
            },
            {
                "title": "Hardcoded Credentials",
                "description": "Hardcoded credentials detected",
                "severity": "critical",
                "cwe": "CWE-798",
                "remediation": "Use environment variables or a secure credential store"
            },
            {
                "title": "Insecure Deserialization",
                "description": "Potential insecure deserialization vulnerability",
                "severity": "high",
                "cwe": "CWE-502",
                "remediation": "Validate and sanitize input before deserialization"
            },
            {
                "title": "Command Injection",
                "description": "Potential command injection vulnerability",
                "severity": "high",
                "cwe": "CWE-78",
                "remediation": "Use safe APIs or properly escape inputs"
            }
        ]

        for file in random.sample(files, min(len(files), 3)):
            for _ in range(random.randint(0, 2)):
                issue = random.choice(issues)
                finding = issue.copy()
                finding["file"] = file
                finding["line"] = str(random.randint(10, 500))
                finding["code"] = f"Example vulnerable code in {file}"
                findings.append(finding)

        return findings

    def _simulate_bandit_findings(self, files: List[str]) -> List[Dict[str, Any]]:
        """Simulate bandit findings for Python files."""
        findings = []

        issues = [
            {
                "title": "Use of insecure MD5 hash function",
                "description": "MD5 is a cryptographically broken hash function",
                "severity": "medium",
                "cwe": "CWE-327",
                "remediation": "Use a secure hashing function like SHA-256"
            },
            {
                "title": "Use of eval()",
                "description": "Use of eval() is insecure",
                "severity": "high",
                "cwe": "CWE-95",
                "remediation": "Avoid using eval() with untrusted input"
            },
            {
                "title": "Possible shell injection",
                "description": "Possible shell injection via subprocess call",
                "severity": "high",
                "cwe": "CWE-78",
                "remediation": "Use subprocess.run with shell=False"
            },
            {
                "title": "Possible SQL injection",
                "description": "SQL injection via string formatting",
                "severity": "high",
                "cwe": "CWE-89",
                "remediation": "Use parameterized queries"
            },
            {
                "title": "Weak cryptography",
                "description": "Use of weak cryptographic algorithm",
                "severity": "medium",
                "cwe": "CWE-326",
                "remediation": "Use strong cryptographic algorithms"
            }
        ]

        for file in random.sample(files, min(len(files), 2)):
            for _ in range(random.randint(0, 2)):
                issue = random.choice(issues)
                finding = issue.copy()
                finding["file"] = file
                finding["line"] = str(random.randint(10, 500))
                finding["code"] = f"Example vulnerable code in {file}"
                findings.append(finding)

        return findings

    def _simulate_gosec_findings(self, files: List[str]) -> List[Dict[str, Any]]:
        """Simulate gosec findings for Go files."""
        return []

    def _simulate_eslint_findings(self, files: List[str]) -> List[Dict[str, Any]]:
        """Simulate eslint-security findings for JavaScript files."""
        return []

    def _simulate_cppcheck_findings(self, files: List[str]) -> List[Dict[str, Any]]:
        """Simulate cppcheck findings for C/C++ files."""
        return []

    def _simulate_sonarqube_findings(self, files: List[str]) -> List[Dict[str, Any]]:
        """Simulate sonarqube findings."""
        return []

    def _simulate_flawfinder_findings(self, files: List[str]) -> List[Dict[str, Any]]:
        """Simulate flawfinder findings for C/C++ files."""
        return []

    def _simulate_brakeman_findings(self, files: List[str]) -> List[Dict[str, Any]]:
        """Simulate brakeman findings for Ruby files."""
        return []

    def _deduplicate_findings(self, findings: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
        """Deduplicate findings based on file, line, and title."""
        unique_findings = {}

        for finding in findings:
            key = (
                finding.get("file", ""),
                finding.get("line", 0),
                finding.get("title", "")
            )

            if key not in unique_findings:
                unique_findings[key] = finding

        return list(unique_findings.values())

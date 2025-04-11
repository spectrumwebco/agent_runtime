import pytest
from pathlib import Path
from unittest.mock import patch, MagicMock, AsyncMock

from agent.environment.repo import (
    LocalRepoConfig,
    GithubRepoConfig,
    GiteeRepoConfig,
    PreExistingRepoConfig,
    repo_from_simplified_input,
)
from agent.utils.gitee import GiteeClient, InvalidGiteeURL


from swerex.deployment.abstract import AbstractDeployment
from swerex.runtime.abstract import CommandResponse

mock_deployment = MagicMock(spec=AbstractDeployment)
mock_deployment.runtime = AsyncMock()
mock_deployment.runtime.execute = AsyncMock(return_value=CommandResponse(exit_code=0, stdout="", stderr=""))
mock_deployment.runtime.upload = AsyncMock()


@pytest.fixture
def mock_git_repo(tmp_path):
    """Fixture to create a dummy git repository."""
    repo_path = tmp_path / "test_local_repo"
    repo_path.mkdir()
    import subprocess
    subprocess.run(["git", "init"], cwd=repo_path, check=True, capture_output=True)
    (repo_path / "test.txt").write_text("hello")
    subprocess.run(["git", "add", "."], cwd=repo_path, check=True, capture_output=True)
    subprocess.run(["git", "commit", "-m", "Initial commit"], cwd=repo_path, check=True, capture_output=True)
    return repo_path

def test_local_repo_config(mock_git_repo):
    config = LocalRepoConfig(path=mock_git_repo)
    assert config.repo_name == "test_local_repo"
    assert config.type == "local"
    assert config.base_commit == "HEAD"
    assert config.get_reset_commands() == [
        "git status",
        "git restore .",
        "git reset --hard HEAD",
        "git clean -fdq",
    ]
    mock_deployment.reset_mock() # Reset mocks for clean assertion
    mock_deployment.runtime.upload = AsyncMock()
    mock_deployment.runtime.execute = AsyncMock(return_value=CommandResponse(exit_code=0, stdout="", stderr=""))

    config.copy(mock_deployment)

    mock_deployment.runtime.upload.assert_awaited_once()
    mock_deployment.runtime.execute.assert_awaited_once() # Should now be awaited correctly for chown

def test_github_repo_config():
    url = "https://github.com/owner/repo"
    config = GithubRepoConfig(github_url=url, base_commit="dev")
    assert config.repo_name == "owner__repo"
    assert config.type == "github"
    assert config.base_commit == "dev"
    assert config.get_reset_commands() == [
        "git status",
        "git restore .",
        "git reset --hard dev",
        "git clean -fdq",
    ]
    mock_deployment.reset_mock() # Reset mocks for clean assertion
    mock_deployment.runtime.execute = AsyncMock(return_value=CommandResponse(exit_code=0, stdout="", stderr=""))
    with patch("os.getenv", return_value="fake_token"):
        config.copy(mock_deployment)
        mock_deployment.runtime.execute.assert_awaited_once() # Use assert_awaited_once
        args, kwargs = mock_deployment.runtime.execute.call_args
        from swerex.runtime.abstract import Command
        command_obj = args[0] # The argument is a Command object
        assert isinstance(command_obj, Command)
        command_str = command_obj.command
        assert "git remote add origin https://fake_token@github.com/owner/repo" in command_str
        assert "git fetch --depth 1 origin dev" in command_str

def test_gitee_repo_config():
    url = "https://gitee.com/owner/repo"
    config = GiteeRepoConfig(gitee_url=url, base_commit="v1.0")
    assert config.repo_name == "owner__repo"
    assert config.type == "gitee"
    assert config.base_commit == "v1.0"
    assert config.get_reset_commands() == [
        "git status",
        "git restore .",
        "git reset --hard v1.0",
        "git clean -fdq",
    ]
    mock_deployment.reset_mock() # Reset mocks for clean assertion
    mock_deployment.runtime.execute = AsyncMock(return_value=CommandResponse(exit_code=0, stdout="", stderr=""))
    with patch("os.getenv", return_value="fake_gitee_token"): # Mock GITEE_ACCESS_TOKEN
        config.copy(mock_deployment)
        mock_deployment.runtime.execute.assert_awaited_once() # Use assert_awaited_once
        args, kwargs = mock_deployment.runtime.execute.call_args
        from swerex.runtime.abstract import Command
        command_obj = args[0] # The argument is a Command object
        assert isinstance(command_obj, Command)
        command_str = command_obj.command
        assert "git remote add origin https://fake_gitee_token@gitee.com/owner/repo" in command_str
        assert "git fetch --depth 1 origin v1.0" in command_str


def test_gitee_parse_url():
    """Test Gitee URL parsing."""
    assert GiteeClient.parse_gitee_url("https://gitee.com/owner/repo") == {"owner": "owner", "repo": "repo"}
    assert GiteeClient.parse_gitee_url("https://gitee.com/owner/repo.git") == {"owner": "owner", "repo": "repo"}
    assert GiteeClient.parse_gitee_url("https://gitee.com/org-name/repo-name") == {"owner": "org-name", "repo": "repo-name"}
    
    with pytest.raises(InvalidGiteeURL):
        GiteeClient.parse_gitee_url("https://github.com/owner/repo")
    with pytest.raises(InvalidGiteeURL):
        GiteeClient.parse_gitee_url("https://gitee.com/owner") # Missing repo
    with pytest.raises(InvalidGiteeURL):
        GiteeClient.parse_gitee_url("invalid-url")


def test_preexisting_repo_config():
    mock_deployment.reset_mock() # Reset mocks to avoid state leakage

    config = PreExistingRepoConfig(repo_name="existing_repo", base_commit="feature-branch")
    assert config.repo_name == "existing_repo"
    assert config.type == "preexisting"
    assert config.base_commit == "feature-branch"
    assert config.get_reset_commands() == [
        "git status",
        "git restore .",
        "git reset --hard feature-branch",
        "git clean -fdq",
    ]
    config.copy(mock_deployment)
    mock_deployment.runtime.upload.assert_not_called()
    mock_deployment.runtime.execute.assert_not_called()


@pytest.mark.parametrize(
    "input_str, expected_type, expected_url_attr, expected_path_attr, base_commit",
    [
        ("https://github.com/owner/repo", GithubRepoConfig, "github_url", None, "main"),
        ("owner/repo", GithubRepoConfig, "github_url", None, "dev"), # Shorthand
        ("https://gitee.com/owner/repo.git", GiteeRepoConfig, "gitee_url", None, "v1.2"),
        ("/local/path/to/repo", LocalRepoConfig, None, "path", "HEAD"),
        ("./relative/path", LocalRepoConfig, None, "path", "tag1"),
    ],
)
@patch('agent.environment.repo_provider.Path') # Patch Path where it's used in repo_provider
def test_repo_from_simplified_input_auto(MockPath, input_str, expected_type, expected_url_attr, expected_path_attr, base_commit, tmp_path):
    original_input_str = input_str

    mock_path_instance = MockPath.return_value
    def exists_side_effect():
        path_str = str(mock_path_instance)
        if path_str.startswith(str(tmp_path)):
            return True
        elif original_input_str.startswith(("/", ".")):
            return True
        return False
    mock_path_instance.exists.side_effect = exists_side_effect

    if expected_path_attr:
        local_path = tmp_path / "test_repos"
        local_path.mkdir(parents=True, exist_ok=True)
        input_str_for_function = str(local_path)
    else:
        input_str_for_function = input_str

    config = repo_from_simplified_input(input=input_str_for_function, base_commit=base_commit, type="auto")

    assert isinstance(config, expected_type)
    assert config.base_commit == base_commit
    if expected_url_attr:
        current_input = original_input_str
        if expected_type == GithubRepoConfig:
            expected_url = f"https://github.com/{current_input}" if '/' in current_input and not current_input.startswith(('http', '/', '.')) else current_input
        elif expected_type == GiteeRepoConfig:
            expected_url = current_input # Assume Gitee URLs are always full in test data
        else: # LocalRepoConfig or others
             expected_url = current_input # Should not have url attr

        if hasattr(config, expected_url_attr):
             assert getattr(config, expected_url_attr) == expected_url
        else:
             assert expected_url_attr is None

    if expected_path_attr:
        assert getattr(config, expected_path_attr) == Path(input_str_for_function)


def test_repo_from_simplified_input_explicit_type(tmp_path):
    github_conf = repo_from_simplified_input(input="owner/repo", type="github")
    assert isinstance(github_conf, GithubRepoConfig)
    assert github_conf.github_url == "https://github.com/owner/repo" # Auto-prefix

    gitee_conf = repo_from_simplified_input(input="https://gitee.com/owner/repo", type="gitee")
    assert isinstance(gitee_conf, GiteeRepoConfig)

    local_path = tmp_path / "local_repo_explicit"
    local_path.mkdir()
    local_conf = repo_from_simplified_input(input=str(local_path), type="local")
    assert isinstance(local_conf, LocalRepoConfig)

    preexisting_conf = repo_from_simplified_input(input="preexisting_name", type="preexisting")
    assert isinstance(preexisting_conf, PreExistingRepoConfig)


def test_repo_from_simplified_input_invalid_auto():
    with pytest.raises(ValueError, match="Could not automatically determine repository type"):
        repo_from_simplified_input(input="invalid-input-format", type="auto")

    with pytest.raises(ValueError, match="Could not automatically determine repository type"):
        repo_from_simplified_input(input="/non/existent/path", type="auto")

def test_repo_from_simplified_input_invalid_type():
     with pytest.raises(ValueError, match="Unknown repo type"):
        repo_from_simplified_input(input="any", type="invalid_type")

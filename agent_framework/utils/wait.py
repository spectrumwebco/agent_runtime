import time
from collections.abc import Callable


async def _wait_until_alive(
    function: Callable, timeout: float = 10.0, function_timeout: float | None = 0.1, sleep: float = 0.25
):
    """Wait until the function returns a truthy value.

    Args:
        function: The function to wait for.
        timeout: The maximum time to wait.
        function_timeout: The timeout passed to the function.
        sleep: The time to sleep between attempts.

    Raises:
        TimeoutError
    """
    end_time = time.time() + timeout
    n_attempts = 0
    await_response = None
    while time.time() < end_time:
        await_response = await function(timeout=function_timeout)
        if await_response:
            return
        time.sleep(sleep)
        n_attempts += 1
    last_response_message = await_response.message if await_response else None
    msg = (
        f"Runtime did not start within {timeout}s (tried to connect {n_attempts} times). "
        f"The last await response was:\n{last_response_message}"
    )
    raise TimeoutError(msg)

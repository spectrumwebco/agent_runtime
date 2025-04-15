from django.apps import AppConfig
import logging

logger = logging.getLogger(__name__)


class PythonMLConfig(AppConfig):
    default_auto_field = "django.db.models.BigAutoField"
    name = "apps.python_ml"
    verbose_name = "Python ML"
    
    def ready(self):
        """
        Initialize the app when Django starts.
        """
        import asyncio
        from .integration.eventstream_integration import event_stream
        
        logger.info("Initializing ML app eventstream integration")
        
        async def start_event_stream():
            try:
                await event_stream.start()
                logger.info("ML app eventstream started successfully")
            except Exception as e:
                logger.error(f"Failed to start eventstream: {e}")
                
        try:
            loop = asyncio.get_event_loop()
            if loop.is_running():
                loop.create_task(start_event_stream())
            else:
                loop.run_until_complete(start_event_stream())
        except Exception as e:
            logger.error(f"Error setting up eventstream: {e}")

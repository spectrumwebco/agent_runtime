"""
Views for billing and payment processing with Polar.sh.

This module provides API views for billing and payment processing
using Polar.sh as the payment processor.
"""

import json
import logging
import requests
from typing import Dict, Any, Optional
from django.http import JsonResponse, HttpResponse
from django.views.decorators.http import require_http_methods
from django.views.decorators.csrf import csrf_exempt
from django.contrib.auth.decorators import login_required
from django.conf import settings
from django.shortcuts import redirect
from rest_framework.decorators import api_view, permission_classes
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

from api.models import UserCredits, CreditTransaction, PolarSubscription

logger = logging.getLogger(__name__)


@api_view(['GET'])
@permission_classes([IsAuthenticated])
def get_credits(request):
    """Get the user's credit balance."""
    credits, created = UserCredits.objects.get_or_create(
        user=request.user,
        defaults={'balance': 0}
    )
    
    transactions = CreditTransaction.objects.filter(
        user=request.user
    ).order_by('-created_at')[:10]
    
    transaction_data = []
    for transaction in transactions:
        transaction_data.append({
            'id': str(transaction.id),
            'amount': float(transaction.amount),
            'transaction_type': transaction.transaction_type,
            'description': transaction.description,
            'created_at': transaction.created_at.isoformat(),
        })
    
    return Response({
        'balance': float(credits.balance),
        'transactions': transaction_data,
    })


@api_view(['POST'])
@permission_classes([IsAuthenticated])
def create_checkout_session(request):
    """Create a Polar.sh checkout session."""
    try:
        plan_id = request.data.get('plan_id')
        
        if not plan_id:
            return Response({
                'status': 'error',
                'message': 'No plan ID provided'
            }, status=400)
        
        response = requests.post(
            f"{settings.POLAR_API_URL}/checkout/session",
            headers={
                'Authorization': f"Bearer {settings.POLAR_API_KEY}",
                'Content-Type': 'application/json',
            },
            json={
                'plan_id': plan_id,
                'customer_email': request.user.email,
                'success_url': f"{settings.FRONTEND_URL}/billing/success?session_id={{CHECKOUT_SESSION_ID}}",
                'cancel_url': f"{settings.FRONTEND_URL}/billing/cancel",
                'metadata': {
                    'user_id': str(request.user.id),
                },
            }
        )
        
        if response.status_code != 200:
            logger.error(f"Error creating Polar.sh checkout session: {response.text}")
            return Response({
                'status': 'error',
                'message': 'Error creating checkout session'
            }, status=500)
        
        checkout_data = response.json()
        checkout_url = checkout_data.get('url')
        
        return Response({
            'status': 'success',
            'checkout_url': checkout_url,
        })
    
    except Exception as e:
        logger.error(f"Error creating checkout session: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@api_view(['POST'])
@permission_classes([IsAuthenticated])
def create_customer_portal_session(request):
    """Create a Polar.sh customer portal session."""
    try:
        subscription = PolarSubscription.objects.filter(user=request.user).first()
        
        if not subscription:
            return Response({
                'status': 'error',
                'message': 'No subscription found'
            }, status=404)
        
        response = requests.post(
            f"{settings.POLAR_API_URL}/customer/portal",
            headers={
                'Authorization': f"Bearer {settings.POLAR_API_KEY}",
                'Content-Type': 'application/json',
            },
            json={
                'subscription_id': subscription.subscription_id,
                'return_url': f"{settings.FRONTEND_URL}/billing",
            }
        )
        
        if response.status_code != 200:
            logger.error(f"Error creating Polar.sh customer portal session: {response.text}")
            return Response({
                'status': 'error',
                'message': 'Error creating customer portal session'
            }, status=500)
        
        portal_data = response.json()
        portal_url = portal_data.get('url')
        
        return Response({
            'status': 'success',
            'portal_url': portal_url,
        })
    
    except Exception as e:
        logger.error(f"Error creating customer portal session: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@csrf_exempt
@require_http_methods(["POST"])
def webhook(request):
    """Handle Polar.sh webhooks."""
    try:
        signature = request.headers.get('Polar-Signature')
        
        if not signature:
            logger.error("No Polar-Signature header found")
            return HttpResponse(status=400)
        
        payload = json.loads(request.body)
        event_type = payload.get('type')
        
        if event_type == 'subscription.created':
            handle_subscription_created(payload)
        elif event_type == 'subscription.updated':
            handle_subscription_updated(payload)
        elif event_type == 'subscription.deleted':
            handle_subscription_deleted(payload)
        elif event_type == 'subscription.payment_succeeded':
            handle_payment_succeeded(payload)
        elif event_type == 'subscription.payment_failed':
            handle_payment_failed(payload)
        
        return HttpResponse(status=200)
    
    except Exception as e:
        logger.error(f"Error handling webhook: {e}")
        return HttpResponse(status=500)


def handle_subscription_created(payload: Dict[str, Any]):
    """Handle subscription.created webhook event."""
    try:
        subscription_data = payload.get('data', {}).get('subscription', {})
        
        if not subscription_data:
            logger.error("No subscription data found in webhook payload")
            return
        
        metadata = subscription_data.get('metadata', {})
        user_id = metadata.get('user_id')
        
        if not user_id:
            logger.error("No user ID found in subscription metadata")
            return
        
        from django.contrib.auth.models import User
        user = User.objects.filter(id=user_id).first()
        
        if not user:
            logger.error(f"User not found: {user_id}")
            return
        
        subscription, created = PolarSubscription.objects.update_or_create(
            user=user,
            defaults={
                'subscription_id': subscription_data.get('id'),
                'plan_id': subscription_data.get('plan_id'),
                'status': subscription_data.get('status'),
                'current_period_start': subscription_data.get('current_period_start'),
                'current_period_end': subscription_data.get('current_period_end'),
                'cancel_at_period_end': subscription_data.get('cancel_at_period_end', False),
            }
        )
        
        plan_id = subscription_data.get('plan_id')
        
        if plan_id == 'basic':
            add_credits(user, 100, 'subscription', 'Basic plan subscription')
        elif plan_id == 'pro':
            add_credits(user, 500, 'subscription', 'Pro plan subscription')
        elif plan_id == 'enterprise':
            add_credits(user, 2000, 'subscription', 'Enterprise plan subscription')
        
        logger.info(f"Subscription created for user {user.username}")
    
    except Exception as e:
        logger.error(f"Error handling subscription.created webhook: {e}")


def handle_subscription_updated(payload: Dict[str, Any]):
    """Handle subscription.updated webhook event."""
    try:
        subscription_data = payload.get('data', {}).get('subscription', {})
        
        if not subscription_data:
            logger.error("No subscription data found in webhook payload")
            return
        
        subscription_id = subscription_data.get('id')
        
        if not subscription_id:
            logger.error("No subscription ID found in webhook payload")
            return
        
        subscription = PolarSubscription.objects.filter(subscription_id=subscription_id).first()
        
        if not subscription:
            logger.error(f"Subscription not found: {subscription_id}")
            return
        
        subscription.plan_id = subscription_data.get('plan_id')
        subscription.status = subscription_data.get('status')
        subscription.current_period_start = subscription_data.get('current_period_start')
        subscription.current_period_end = subscription_data.get('current_period_end')
        subscription.cancel_at_period_end = subscription_data.get('cancel_at_period_end', False)
        subscription.save()
        
        logger.info(f"Subscription updated for user {subscription.user.username}")
    
    except Exception as e:
        logger.error(f"Error handling subscription.updated webhook: {e}")


def handle_subscription_deleted(payload: Dict[str, Any]):
    """Handle subscription.deleted webhook event."""
    try:
        subscription_data = payload.get('data', {}).get('subscription', {})
        
        if not subscription_data:
            logger.error("No subscription data found in webhook payload")
            return
        
        subscription_id = subscription_data.get('id')
        
        if not subscription_id:
            logger.error("No subscription ID found in webhook payload")
            return
        
        subscription = PolarSubscription.objects.filter(subscription_id=subscription_id).first()
        
        if not subscription:
            logger.error(f"Subscription not found: {subscription_id}")
            return
        
        subscription.delete()
        
        logger.info(f"Subscription deleted for user {subscription.user.username}")
    
    except Exception as e:
        logger.error(f"Error handling subscription.deleted webhook: {e}")


def handle_payment_succeeded(payload: Dict[str, Any]):
    """Handle subscription.payment_succeeded webhook event."""
    try:
        subscription_data = payload.get('data', {}).get('subscription', {})
        
        if not subscription_data:
            logger.error("No subscription data found in webhook payload")
            return
        
        subscription_id = subscription_data.get('id')
        
        if not subscription_id:
            logger.error("No subscription ID found in webhook payload")
            return
        
        subscription = PolarSubscription.objects.filter(subscription_id=subscription_id).first()
        
        if not subscription:
            logger.error(f"Subscription not found: {subscription_id}")
            return
        
        plan_id = subscription_data.get('plan_id')
        
        if plan_id == 'basic':
            add_credits(subscription.user, 100, 'subscription', 'Basic plan subscription renewal')
        elif plan_id == 'pro':
            add_credits(subscription.user, 500, 'subscription', 'Pro plan subscription renewal')
        elif plan_id == 'enterprise':
            add_credits(subscription.user, 2000, 'subscription', 'Enterprise plan subscription renewal')
        
        logger.info(f"Payment succeeded for user {subscription.user.username}")
    
    except Exception as e:
        logger.error(f"Error handling subscription.payment_succeeded webhook: {e}")


def handle_payment_failed(payload: Dict[str, Any]):
    """Handle subscription.payment_failed webhook event."""
    try:
        subscription_data = payload.get('data', {}).get('subscription', {})
        
        if not subscription_data:
            logger.error("No subscription data found in webhook payload")
            return
        
        subscription_id = subscription_data.get('id')
        
        if not subscription_id:
            logger.error("No subscription ID found in webhook payload")
            return
        
        subscription = PolarSubscription.objects.filter(subscription_id=subscription_id).first()
        
        if not subscription:
            logger.error(f"Subscription not found: {subscription_id}")
            return
        
        subscription.status = 'past_due'
        subscription.save()
        
        logger.info(f"Payment failed for user {subscription.user.username}")
    
    except Exception as e:
        logger.error(f"Error handling subscription.payment_failed webhook: {e}")


def add_credits(user, amount, transaction_type, description):
    """Add credits to a user's account."""
    try:
        credits, created = UserCredits.objects.get_or_create(
            user=user,
            defaults={'balance': 0}
        )
        
        credits.balance += amount
        credits.save()
        
        CreditTransaction.objects.create(
            user=user,
            amount=amount,
            transaction_type=transaction_type,
            description=description
        )
        
        logger.info(f"Added {amount} credits to user {user.username}")
        
        return True
    
    except Exception as e:
        logger.error(f"Error adding credits: {e}")
        return False

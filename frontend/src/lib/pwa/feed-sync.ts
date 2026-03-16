import { API_BASE_URL } from '$lib/config/env';
import { enqueueRequest } from '$lib/pwa/offline-queue';
import { scheduleOutboxSync } from '$lib/pwa/runtime';

export async function queueCreateActivity(payload: {
  type: string;
  title: string;
  content: string;
  points?: number;
}) {
  await enqueueRequest({
    type: 'create-activity',
    url: `${API_BASE_URL}/activities`,
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(payload)
  });

  await scheduleOutboxSync();
}

export async function queueCreateComment(payload: {
  activityID: string;
  comment: string;
}) {
  await enqueueRequest({
    type: 'create-comment',
    url: `${API_BASE_URL}/activities/${payload.activityID}/comments`,
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      comment: payload.comment
    })
  });

  await scheduleOutboxSync();
}

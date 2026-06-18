import { useAuthStore } from '../stores/authStore'

type HTTPMethod = 'GET' | 'POST' | 'PUT' | 'DELETE'

export async function protectedFetch(url: string, method: HTTPMethod, data?: unknown) {
    const token = useAuthStore.getState().token

    const options: RequestInit = {
        method,
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
        },
    }

    if (data) options.body = JSON.stringify(data)

    return fetch(`http://localhost:8080${url}`, options)
}

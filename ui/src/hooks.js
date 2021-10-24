import * as cookie from "cookie"

export function getSession(request) {
	const cookies = cookie.parse(request.headers.cookie || "")
	return {
		access_token: cookies.access_token,
		expires_at: cookies.expires_at,
	}
}

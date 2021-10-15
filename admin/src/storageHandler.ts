export default {
    sessionId() {
        let cookie = document.cookie.split('; ').find(row => row.startsWith('session_id'));
        cookie = cookie ? cookie.split('=')[1] : "null";
        return cookie
    }
}
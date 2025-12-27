// refresh.js
async function refreshAccessToken() {
    try {
      const res = await fetch('/refresh-token', {
        method: 'POST',
        credentials: 'include', // важно для работы с cookie
      });
  
      if (!res.ok) {
        console.log('Refresh token failed, user needs to login');
        return false;
      }
  
      const data = await res.json();
      console.log(data.status); // "access token refreshed"
      return true;
    } catch (err) {
      console.error('Error refreshing token:', err);
      return false;
    }
  }
  
  // Можно вызывать на каждой странице при загрузке
  window.addEventListener('load', () => {
    refreshAccessToken();
  });
  
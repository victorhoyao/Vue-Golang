import { setToken } from './auth';

// Test login function to set a valid token for testing
const testLogin = async () => {
  try {
    const response = await fetch('http://localhost:8036/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        userName: 'testadmin',
        passWord: 'test123',
      }),
    });

    const data = await response.json();

    if (data.code === 200 && data.data.token) {
      setToken(data.data.token);
      console.log('Test login successful, token set:', data.data.token);
      return true;
    }
    console.error('Test login failed:', data);
    return false;
  } catch (error) {
    console.error('Test login error:', error);
    return false;
  }
};

// Function to call from browser console for testing
(window as any).testLogin = testLogin;

export default testLogin; 
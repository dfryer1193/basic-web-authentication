const registerForm = document.getElementById('registerForm');
const loginForm = document.getElementById('loginForm');
const changePasswordForm = document.getElementById('changePasswordForm');
const changePasswordMessage = document.getElementById('changePasswordMessage');
const welcomeBtn = document.getElementById('welcomeBtn');
const welcomeMessage = document.getElementById('welcomeMessage');

registerForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  const username = document.getElementById('regUsername').value;
  const password = document.getElementById('regPassword').value;

  const res = await fetch('/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });

  if (res.ok) {
    alert('User registered successfully');
  } else {
    alert('Error registering user');
  }
});

loginForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  const username = document.getElementById('loginUsername').value;
  const password = document.getElementById('loginPassword').value;

  const res = await fetch('/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });

  if (res.ok) {
    alert('Login successful');
  } else {
    alert('Invalid username or password');
  }
});

changePasswordForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  const newPassword = document.getElementById('newPassword').value;

  const res = await fetch('/change-password', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ newPassword }),
  });

  if (res.ok) {
    changePasswordMessage.textContent = 'Password changed successfully!';
    changePasswordMessage.style.color = 'green';
  } else {
    const errorMessage = await res.text();
    changePasswordMessage.textContent = `Error: ${errorMessage}`;
    changePasswordMessage.style.color = 'red';
  }
});

welcomeBtn.addEventListener('click', async () => {
  const res = await fetch('/welcome', { method: 'GET' });

  if (res.ok) {
    welcomeMessage.textContent = await res.text();
  } else {
    alert('You are not authorized');
  }
});
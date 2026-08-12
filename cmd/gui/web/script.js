// --- Background theme song ---
// Most browsers block autoplay with sound until the user has interacted
// with the page at least once. Retry playback on the first click/keypress
// if the initial autoplay attempt was blocked.
const theme = document.getElementById('theme');
theme.play().catch(() => {
  const resume = () => {
    theme.play();
    document.removeEventListener('click', resume);
    document.removeEventListener('keydown', resume);
  };
  document.addEventListener('click', resume);
  document.addEventListener('keydown', resume);
});

// --- UI sound effects ---
const SFX = {
  hover: '/assets/audios/Voicy_%20Item%20Get.mp3',
  error: '/assets/audios/Voicy_Error%20Select.mp3',
  confirm: '/assets/audios/Voicy_Select.mp3',
};

function playSound(path) {
  new Audio(path).play().catch(() => {});
}

const usernameInput = document.getElementById('username');
const connectButton = document.querySelector('#connect-form button');
const status = document.getElementById('status');

[usernameInput, connectButton].forEach((el) => {
  el.addEventListener('focus', () => playSound(SFX.hover));
});

document.getElementById('connect-form').addEventListener('submit', (event) => {
  event.preventDefault();

  if (usernameInput.value.trim() === '') {
    playSound(SFX.error);
    status.textContent = 'Username is required.';
    status.classList.add('error');
    return;
  }

  status.classList.remove('error');
  playSound(SFX.confirm);
  status.textContent = 'Connecting...';
  // TODO: open the WebSocket connection here.
});

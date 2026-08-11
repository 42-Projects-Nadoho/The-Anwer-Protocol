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

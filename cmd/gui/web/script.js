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

// --- Elements ---
const loginView = document.getElementById('login-view');
const gameView = document.getElementById('game-view');
const usernameInput = document.getElementById('username');
const connectButton = document.querySelector('#connect-form button');
const status = document.getElementById('status');

const roomName = document.getElementById('room-name');
const roomDesc = document.getElementById('room-desc');
const roomPlayers = document.getElementById('room-players');
const exitsBox = document.getElementById('exits');
const whoCount = document.getElementById('who-count');
const groupStatus = document.getElementById('group-status');
const chatLog = document.getElementById('chat-log');
const eventLog = document.getElementById('event-log');

[usernameInput, connectButton].forEach((el) => {
  el.addEventListener('focus', () => playSound(SFX.hover));
});

// --- Tabs ---
document.querySelectorAll('.tab-btn').forEach((btn) => {
  btn.addEventListener('click', () => {
    document.querySelectorAll('.tab-btn').forEach((b) => b.classList.remove('active'));
    document.querySelectorAll('.tab-panel').forEach((p) => p.classList.remove('active'));
    btn.classList.add('active');
    document.getElementById(btn.dataset.tab + '-tab').classList.add('active');
  });
});

// --- Connection ---
let ws = null;
let myGroupID = '';

document.getElementById('connect-form').addEventListener('submit', (event) => {
  event.preventDefault();

  const username = usernameInput.value.trim();
  if (username === '') {
    playSound(SFX.error);
    status.textContent = 'Username is required.';
    status.classList.add('error');
    return;
  }

  status.classList.remove('error');
  playSound(SFX.confirm);
  status.textContent = 'Connecting...';
  connectButton.disabled = true;
  openConnection(username);
});

function openConnection(username) {
  ws = new WebSocket('ws://' + location.host + '/ws');

  ws.onerror = () => {
    status.textContent = 'Connection failed.';
    status.classList.add('error');
    connectButton.disabled = false;
  };

  ws.onclose = () => {
    if (!gameView.classList.contains('hidden')) {
      appendLog('--- disconnected ---');
    } else {
      connectButton.disabled = false;
    }
  };

  ws.onmessage = (event) => handleLine(event.data, username);
}

function send(cmd) {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(cmd);
  }
}

// --- Server line dispatch (mirrors the CLI's own logic) ---
function handleLine(line, username) {
  if (line === 'OK hello proto=1') {
    send('CONNECT ' + username);
    return;
  }

  if (line === 'OK connected') {
    showGameView();
    send('LOOK');
    return;
  }

  if (line.startsWith('ERR 201 NAME_IN_USE')) {
    status.textContent = 'Username already taken.';
    status.classList.add('error');
    connectButton.disabled = false;
    ws.close();
    return;
  }

  if (line.startsWith('OK {')) {
    renderRoom(JSON.parse(line.slice(3)));
    return;
  }

  if (line.startsWith('OK room=')) {
    send('LOOK');
    return;
  }

  if (line.startsWith('OK players=')) {
    whoCount.textContent = 'Players online: ' + line.slice('OK players='.length);
    return;
  }

  if (line.startsWith('OK group=')) {
    myGroupID = line.slice('OK group='.length);
    groupStatus.textContent = 'Group: ' + myGroupID;
    appendLog(line);
    return;
  }

  if (line === 'OK bye') {
    ws.close();
    return;
  }

  if (/^EVT (GLOBAL|ROOM|GROUP) CHAT /.test(line)) {
    appendChat(line);
    return;
  }

  if (line.startsWith('EVT GROUP LEAVE') && line.includes(username)) {
    myGroupID = '';
    groupStatus.textContent = 'Not in a group.';
  }

  appendLog(line);
}

function showGameView() {
  loginView.classList.add('hidden');
  gameView.classList.remove('hidden');
}

function renderRoom(data) {
  roomName.textContent = data.room.name;
  roomDesc.textContent = data.room.description;
  roomPlayers.textContent = 'Players here: ' + data.players.join(', ');

  exitsBox.innerHTML = '';
  Object.keys(data.room.exits).forEach((direction) => {
    const btn = document.createElement('button');
    btn.type = 'button';
    btn.textContent = direction;
    btn.addEventListener('click', () => send('MOVE ' + direction));
    exitsBox.appendChild(btn);
  });
}

function appendChat(line) {
  chatLog.textContent += line + '\n';
  chatLog.scrollTop = chatLog.scrollHeight;
}

function appendLog(line) {
  eventLog.textContent += line + '\n';
  eventLog.scrollTop = eventLog.scrollHeight;
}

// --- Chat ---
document.getElementById('chat-form').addEventListener('submit', (event) => {
  event.preventDefault();
  const scope = document.getElementById('chat-scope').value;
  const input = document.getElementById('chat-input');
  const message = input.value.trim();
  if (message === '') {
    return;
  }
  send('CHAT ' + scope + ' ' + message);
  input.value = '';
});

// --- Group actions ---
document.getElementById('group-create-btn').addEventListener('click', () => {
  send('GROUP CREATE');
});

document.getElementById('group-leave-btn').addEventListener('click', () => {
  send('GROUP LEAVE');
  myGroupID = '';
  groupStatus.textContent = 'Not in a group.';
});

document.getElementById('group-invite-form').addEventListener('submit', (event) => {
  event.preventDefault();
  const input = document.getElementById('invite-username');
  const name = input.value.trim();
  if (name === '') {
    return;
  }
  send('GROUP INVITE ' + name);
  input.value = '';
});

document.getElementById('group-join-form').addEventListener('submit', (event) => {
  event.preventDefault();
  const input = document.getElementById('join-groupid');
  const groupID = input.value.trim();
  if (groupID === '') {
    return;
  }
  send('GROUP JOIN ' + groupID);
  input.value = '';
});

// --- Top bar actions ---
document.getElementById('who-btn').addEventListener('click', () => send('WHO'));
document.getElementById('quit-btn').addEventListener('click', () => send('QUIT'));

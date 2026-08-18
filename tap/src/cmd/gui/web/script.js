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
const roomItems = document.getElementById('room-items');
const roomNpcs = document.getElementById('room-npcs');
const playerHp = document.getElementById('player-hp');
const statusBtn = document.getElementById('status-btn');
const inventoryList = document.getElementById('inventory-list');
const invBtn = document.getElementById('inv-btn');
const questsList = document.getElementById('quests-list');
const questsBtn = document.getElementById('quests-btn');

if (statusBtn) statusBtn.addEventListener('click', () => send('STATUS'));
if (invBtn) invBtn.addEventListener('click', () => send('INVENTORY'));
if (questsBtn) questsBtn.addEventListener('click', () => send('QUESTS'));
const lookBtn = document.getElementById('look-btn');
if (lookBtn) lookBtn.addEventListener('click', () => send('LOOK'));
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
    status.textContent = 'ERR 900 CONNECTION_FAILED';
    status.classList.add('error');
    connectButton.disabled = false;
  };

  ws.onclose = () => {
    if (!gameView.classList.contains('hidden')) {
      gameView.classList.add('hidden');
      loginView.classList.remove('hidden');
      status.textContent = 'Disconnected.';
      status.classList.remove('error');
    }
    connectButton.disabled = false;
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
    try {
      const data = JSON.parse(line.slice(3));
      if (data.room) {
        renderRoom(data);
      } else if (data.hp !== undefined) {
        playerHp.textContent = `HP: ${data.hp}/${data.max_hp} (${data.status})`;
      } else if (data.quest_id !== undefined) {
        appendChat(`[System] Quest: ${data.quest_id} - ${data.description}`);
        send('QUESTS');
      } else if (data.npc !== undefined) {
        appendChat(`[${data.npc}]: ${data.dialogue}`);
      } else if (data.attacker_hp !== undefined) {
        appendChat(`[Combat] You hit for ${data.damage}. Target HP: ${data.target_hp}. Your HP: ${data.attacker_hp}. Status: ${data.status}`);
        send('STATUS');
        if (data.status === 'victory') {
          send('LOOK');
          send('QUESTS');
          send('INVENTORY');
        }
      } else {
        appendLog(line);
      }
    } catch(e) {
      appendLog(line);
    }
    return;
  }

  if (line.startsWith('OK [')) {
    try {
      const data = JSON.parse(line.slice(3));
      if (data.length === 0 || typeof data[0] === 'string') {
        renderInventory(data);
      } else {
        renderQuests(data);
      }
    } catch(e) {
      appendLog(line);
    }
    return;
  }

  if (line.startsWith('OK taken=')) {
    appendChat(`[System] Item taken.`);
    send('LOOK');
    send('INVENTORY');
    return;
  }

  if (line.startsWith('OK dropped=')) {
    appendChat(`[System] Item dropped.`);
    send('LOOK');
    send('INVENTORY');
    return;
  }

  if (line.startsWith('OK room=')) {
    send('LOOK');
    return;
  }

  if (line.startsWith('OK players=')) {
    const text = line.slice('OK players='.length);
    whoCount.textContent = 'Players online: ' + text.split(' ')[0]; // Just the count for the top bar
    appendChat(`[System] Players online: ${text}`);
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

  if (/^EVT ROOM (COMBAT|RESPAWN) /.test(line)) {
    appendChat(`[System] ${line.slice(9)}`);
    return;
  }

  if (/^EVT ROOM PRESENCE (ENTER|LEAVE) /.test(line)) {
    appendChat(`[System] ${line.slice(9)}`);
    send('LOOK');
    send('WHO');
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
  send('STATUS');
  send('INVENTORY');
  send('QUESTS');
}

function renderInventory(items) {
  inventoryList.innerHTML = '';
  if (!items || items.length === 0) {
    inventoryList.innerHTML = '<li style="color:#ccc; font-size:0.9em;">Empty</li>';
    return;
  }
  items.forEach(item => {
    const li = document.createElement('li');
    li.style.display = 'flex';
    li.style.justifyContent = 'space-between';
    li.style.alignItems = 'center';
    li.style.marginBottom = '4px';
    
    const span = document.createElement('span');
    span.textContent = item;
    
    const btn = document.createElement('button');
    btn.type = 'button';
    btn.textContent = 'Drop';
    btn.style.padding = '2px 5px';
    btn.addEventListener('click', () => send('DROP ' + item));
    
    li.appendChild(span);
    li.appendChild(btn);
    inventoryList.appendChild(li);
  });
}

function renderQuests(quests) {
  questsList.innerHTML = '';
  if (!quests || quests.length === 0) {
    questsList.innerHTML = '<li style="color:#ccc; font-size:0.9em;">None</li>';
    return;
  }
  quests.forEach(q => {
    const li = document.createElement('li');
    li.style.marginBottom = '4px';
    li.textContent = `[${q.status.toUpperCase()}] ${q.quest_id}`;
    if (q.progress) li.textContent += ` (${q.progress})`;
    questsList.appendChild(li);
  });
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
  
  roomItems.innerHTML = '';
  if (!data.items || data.items.length === 0) {
    roomItems.textContent = 'None';
  } else {
    data.items.forEach(item => {
      const row = document.createElement('div');
      row.className = 'row';
      row.style.alignItems = 'center';
      row.style.marginBottom = '5px';
      
      const span = document.createElement('span');
      span.textContent = item;
      span.style.flex = '1';
      
      const btn = document.createElement('button');
      btn.type = 'button';
      btn.textContent = 'Take';
      btn.style.padding = '2px 5px';
      btn.addEventListener('click', () => send('TAKE ' + item));
      
      row.appendChild(span);
      row.appendChild(btn);
      roomItems.appendChild(row);
    });
  }
  
  roomNpcs.innerHTML = '';
  if (!data.npcs || data.npcs.length === 0) {
    roomNpcs.textContent = 'None';
  } else {
    data.npcs.forEach(npc => {
      const row = document.createElement('div');
      row.className = 'row';
      row.style.alignItems = 'center';
      row.style.marginBottom = '5px';
      
      const span = document.createElement('span');
      span.textContent = npc;
      span.style.flex = '1';
      
      const btnTalk = document.createElement('button');
      btnTalk.type = 'button';
      btnTalk.textContent = 'Talk';
      btnTalk.style.padding = '2px 5px';
      btnTalk.addEventListener('click', () => send('TALK ' + npc));
      
      const btnAttack = document.createElement('button');
      btnAttack.type = 'button';
      btnAttack.textContent = 'Atk';
      btnAttack.style.padding = '2px 5px';
      btnAttack.addEventListener('click', () => send('ATTACK ' + npc));
      
      const btnQuest = document.createElement('button');
      btnQuest.type = 'button';
      btnQuest.textContent = 'Quest';
      btnQuest.style.padding = '2px 5px';
      btnQuest.addEventListener('click', () => send('QUEST ' + npc));
      
      row.appendChild(span);
      row.appendChild(btnTalk);
      row.appendChild(btnAttack);
      row.appendChild(btnQuest);
      roomNpcs.appendChild(row);
    });
  }
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

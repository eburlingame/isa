import { v4 as uuid } from "uuid";

const requestTimeout = 5000;
const sessionIdKey = `isa_game_session_id`;

class GameClient {
  constructor(clientId, url) {
    this.url = url;
    this.clientId = clientId;

    this.stateCallbacks = [];

    this._initialize();
  }

  _initialize() {
    console.log(`Attempting to connect to websocket on ${this.url}`);

    this.ws = new WebSocket(this.url);
    this.ws.onopen = this._onOpen.bind(this);
    this.ws.onmessage = this._onMessage.bind(this);
    this.ws.onclose = this._onClose.bind(this);

    this.inflight = {};
    this.inflightTimeouts = {};

    this.state = {
      open: false,
      gameId: null,
      sessionId: null,
      playerName: null,
      isHost: false,
      gameState: null,
    };
  }

  _loadSessionId() {
    const buffer = localStorage.getItem(sessionIdKey);

    if (buffer) {
      return buffer;
    }

    return false;
  }

  _saveSessionId(sessionId) {
    console.log("Writing session id to local storage: " + sessionId);
    localStorage.setItem(sessionIdKey, sessionId);
  }

  _sendCommand(verb, data = {}) {
    const reqId = uuid();
    const payload = { v: verb, d: data, reqId };

    console.log(payload);
    this.ws.send(JSON.stringify(payload));

    return reqId;
  }

  _enqueueCommand(verb, data = {}) {
    let resolver, rejecter;

    // Send the command
    const reqId = this._sendCommand(verb, data);

    // Construct a new promise, and extract the resolve/reject fns
    const promise = new Promise((resolve, reject) => {
      resolver = resolve;
      rejecter = reject;
    });

    // Cache the resolver and rejecter to join them when we receive the response
    this.inflight[reqId] = {
      resolver,
      rejecter,
    };

    // Set a timeout
    this.inflightTimeouts[reqId] = setTimeout(() => {
      const request = { reqId, v: verb, d: data };

      this._clearInflights(reqId);

      rejecter(
        new Error(`Request timed out. Request: ${JSON.stringify(request)}`)
      );
    }, requestTimeout);

    // return the promise
    return promise;
  }

  _clearInflights(reqId) {
    clearTimeout(this.inflightTimeouts[reqId]);
    delete this.inflight[reqId];
    delete this.inflightTimeouts[reqId];
  }

  _onOpen() {
    const contents = this._loadSessionId();

    if (contents) {
      console.log("Loaded sessionId from localStorage: " + contents);
      return this._enqueueCommand("openSession", { sessionId: contents });
    } else {
      return this._enqueueCommand("openSession", {});
    }
  }

  _resolveResponse(resp) {
    const { reqId } = resp;

    if (reqId in this.inflight) {
      this.inflight[reqId].resolver(resp);
      this._clearInflights(reqId);
    }
  }

  _rejectResponse(resp) {
    const { reqId } = resp;

    if (reqId in this.inflight) {
      this.inflight[reqId].rejecter(resp.d.message);
      this._clearInflights(reqId);
    }
  }

  _onMessage({ data }) {
    let msg;

    try {
      msg = JSON.parse(data);
      console.log(msg);
    } catch (e) {
      console.error("Error parsing incoming message: ", data, e);
      return;
    }

    if (msg.err) {
      console.error("Error message received: " + msg.d.message);
      this._rejectResponse(msg);
      return;
    }

    switch (msg.v) {
      case "openSession":
        this.state.open = true;
        this.state.sessionId = msg.d.sessionId;
        this._saveSessionId(this.state.sessionId);
        this._flushStateChange();

        if (this.openResolver) {
          this.openResolver();
          this.openResolver = undefined;
        }

        this._resolveResponse(msg);
        break;

      case "gameState":
        this._processGameUpdate(msg);
        this._resolveResponse(msg);
        break;

      default:
        this._resolveResponse(msg);
        break;
    }
  }

  _onClose() {
    console.log("Connection closed");
    this._initialize();
    this._flushStateChange();
  }

  _processGameUpdate(response) {
    const { abandoned } = response.d;

    if (abandoned) {
      this.state.gameState = null;
    } else {
      this.state.gameId = response.d.gameId;
      this.state.gamePneumonic = response.d.gamePneumonic;
      this.state.gameState = response.d.game;
      this.state.isHost = response.d.isHost;
    }

    this._flushStateChange();

    return {
      gameId: this.state.gameId,
      game: this.state.gameState,
      isHost: this.state.isHost,
    };
  }

  _flushStateChange() {
    console.log(this.state);

    this.stateCallbacks.map((cb) => {
      if (cb) {
        cb(this.state);
      }
    });
  }

  async waitForConnection() {
    const handler = (resolve, reject) => {
      this.openResolver = resolve;
      this.openrejecter = reject;
    };

    return new Promise(handler.bind(this));
  }

  async waitForState() {
    const handler = (resolve, reject) => {
      this.stateResolver = resolve;
      this.staterejecter = reject;
    };

    return new Promise(handler.bind(this));
  }

  getGameId() {
    return this.state.gameId;
  }

  async createGame(playerName) {
    this.state.playerName = playerName;

    const response = await this._enqueueCommand("createGame", {
      playerName: this.state.playerName,
    });

    return this._processGameUpdate(response);
  }

  async joinGame(playerName, gameId) {
    this.state.playerName = playerName;

    const response = await this._enqueueCommand("joinGame", {
      gameId: gameId.trim().toUpperCase(),
      playerName: this.state.playerName,
    });

    return this._processGameUpdate(response);
  }

  async leaveGame() {
    if (!this.state.gameState) {
      console.warn("Not playing a game");
      return;
    }

    const response = await this._enqueueCommand("leaveGame", {
      gameId: this.state.gameId,
      playerName: this.state.playerName,
    });

    return this._processGameUpdate(response);
  }

  async startGame() {
    const response = await this._enqueueCommand("startGame");
    return this._processGameUpdate(response);
  }

  async endGame() {
    const response = await this._enqueueCommand("endGame");
    return this._processGameUpdate(response);
  }

  async playCard(cardIndex, wildColor) {
    const response = await this._enqueueCommand("playCard", {
      cardIndex: cardIndex.toString(),
      wildColor,
    });
    return this._processGameUpdate(response);
  }

  async drawCard() {
    const response = await this._enqueueCommand("drawCard");
    return this._processGameUpdate(response);
  }

  async doneDrawing() {
    const response = await this._enqueueCommand("doneDrawing");
    return this._processGameUpdate(response);
  }

  getGameState() {
    return this.state.gameState;
  }

  subscribeToStateChanges(callback) {
    this.stateCallbacks.push(callback);
  }

  unsubscribeFromStateChanges(callback) {
    for (let i = 0; i < this.stateCallbacks.length; i++) {
      if (this.stateCallbacks[i] === callback) {
        delete this.stateCallbacks[i];
        break;
      }
    }
  }
}

export default GameClient;

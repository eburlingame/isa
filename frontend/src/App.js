import React, { useEffect, useRef, useState } from "react";
import styled from "styled-components";
import { ClientProvider, useClient } from "./ClientContext";
import { useClientState } from "./hooks";

import Game from "./pages/Game";
import Lobby from "./pages/Lobby";
import Summary from "./pages/Summary";
import Welcome from "./pages/Welcome";

const Container = styled.div`
  height: 100%;
  width: 100%;
  background-color: #fcfcfc;

  position: relative;
`;

const windowHost = window.location.host.split(":")[0];

const isLocalAddress =
  windowHost.startsWith("192") || windowHost.startsWith("localhost");

const server = isLocalAddress
  ? `ws://localhost:5000/ws`
  : `wss://${windowHost.replace("uno", "uno-api")}/ws`;

const Consumer = ({}) => {
  const state = useClientState();

  if (!state || state.open === false) {
    return <div>Unable to connect to server</div>;
  }

  if (!state.gameState) {
    return (
      <Container>
        <Welcome />
      </Container>
    );
  }

  const { gameId, gamePneumonic, gameState: game, isHost } = state;
  const { state: status } = game;

  const commonProps = {
    gameCode: gameId,
    gamePneumonic,
    game,
    isHost,
  };

  return (
    <Container>
      {status === 0 && <Lobby {...commonProps} />}
      {status === 1 && <Game {...commonProps} />}
      {status === 2 && <Summary {...commonProps} />}
    </Container>
  );
};

export default ({}) => {
  return (
    <ClientProvider url={server}>
      <Consumer />
    </ClientProvider>
  );
};

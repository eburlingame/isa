import { first } from "lodash";
import { useEffect, useRef, useState } from "react";
import styled from "styled-components";
import { colors } from "../constants";
import ConfettiGenerator from "confetti-js";

import { useActions } from "../hooks";

const Container = styled.div`
  width: 100%;
  height: 100%;

  font-size: 1.25em;

  display: flex;
  justify-content: center;
  align-items: center;
`;

const BackgroundCanvas = styled.canvas`
  position: absolute;
  width: 100%;
  height: 100%;
  z-index: 0;
`;

const Frame = styled.div`
  z-index: 1;

  width: 35vw;
  max-width: 400px;
  min-width: 300px;

  background-color: ${colors.lightGray};
  text-align: center;
  padding: 1em;
  border-radius: 0.5em;
`;

const Button = styled.div`
  cursor: pointer;
  padding: 0.5em;
  border-radius: 0.5em;

  transition: 0.3s background-color;
  background-color: ${colors.accent};
  color: ${colors.white};

  :hover {
    background-color: ${colors.accentDark};
    transition: 0.3s background-color;
  }
`;

const PlayButton = styled(Button)`
  background-color: ${colors.green};
  margin-bottom: 0.25em;
`;

const PlayerContainer = styled.div`
  border: 0.5px solid ${colors.accentDark};
  border-radius: 0.25em;
  margin-bottom: 0.5em;
  padding: 0.25em;
`;

const Player = ({ name, numCards }) => (
  <PlayerContainer>
    {name} had {numCards} cards
  </PlayerContainer>
);

export default ({ gameCode, game, isHost }) => {
  const { startGame, endGame, leaveGame } = useActions();

  const winner = first(
    game.otherPlayers.filter(({ numCards }) => numCards === 0)
  );

  useEffect(() => {
    const confetti = new ConfettiGenerator({
      target: "background-canvas",
      max: 500,
      clock: 30,
    });
    confetti.render();

    return () => confetti.clear();
  }, []);

  return (
    <Container>
      <BackgroundCanvas id="background-canvas" />

      <Frame>
        <h2>{winner.name} Won!</h2>
        <h3>Players:</h3>
        {game.otherPlayers
          .filter((p) => p !== winner)
          .map(({ name, numCards }) => (
            <Player name={name} numCards={numCards}></Player>
          ))}

        {isHost && game.otherPlayers.length > 1 && (
          <PlayButton onClick={() => startGame()}>Play again!</PlayButton>
        )}

        {isHost && <Button onClick={() => endGame()}>End Game</Button>}
        {!isHost && <Button onClick={() => leaveGame()}>Leave Game</Button>}
      </Frame>
    </Container>
  );
};

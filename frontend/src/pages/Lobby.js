import { useState } from "react";
import styled from "styled-components";
import { colors } from "../constants";

import { useActions } from "../hooks";

const Container = styled.div`
  width: 100%;
  height: 100%;

  font-size: 1.25em;

  display: flex;
  justify-content: center;
  align-items: center;
`;

const Frame = styled.div`
  width: 35vw;
  max-width: 400px;
  min-width: 300px;

  background-color: ${colors.lightGray};
  text-align: center;
  padding: 1em;
  border-radius: 0.5em;
`;

const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

const Input = styled.input`
  font-size: 1.2em;
  border-radius: 0.5em;
  padding: 0.5em;
  border: none;
  margin-bottom: 0.5em;

  :focus {
    outline: none;
  }
`;

const Button = styled.div`
  cursor: pointer;
  padding: 0.5em;
  border-radius: 0.5em;

  margin-top: 0.25em;

  transition: 0.3s background-color;
  background-color: ${colors.accent};
  color: ${colors.white};

  :hover {
    background-color: ${colors.accentDark};
    transition: 0.3s background-color;
  }
`;

const PositiveButton = styled(Button)`
  background-color: ${colors.green};
  :hover {
    background-color: ${colors.greenDark}
  }
`

const NegativeButton = styled(Button)`
  background-color: ${colors.red};
  :hover {
    background-color: ${colors.redDark}
  }
`

const PlayerContainer = styled.div`
  border: 0.5px solid ${colors.accentDark};
  border-radius: 0.25em;
  margin-bottom: 0.5em;
  padding: 0.25em;
`;

const Label = styled.div`
  font-size: 1.25em;
`;

const Title = styled.div`
  font-weight: 600;
  font-size: 2em;
`;

const Subtitle = styled.div`
  margin-top: 1.25em;
  font-size: 1.0em;
`

const Pneumonic = styled.div`
  margin-bottom: 1.25em;
  font-size: 1.3em;
`;

const Player = ({ name }) => <PlayerContainer>{name}</PlayerContainer>;

export default ({ gameCode, gamePneumonic, game, isHost }) => {
  const { startGame, endGame, leaveGame } = useActions();

  return (
    <Container>
      <Frame>
        <Label>Get Ready!</Label>

        <Subtitle>Game Code:</Subtitle>
        <Title>{gameCode}</Title>
        <Pneumonic>{gamePneumonic}</Pneumonic>

        <h3>Players:</h3>
        {game.otherPlayers.map(({ name }) => (
          <Player name={name} key={name} />
        ))}

        {isHost && game.otherPlayers.length > 1 && (
          <PositiveButton onClick={() => startGame()}>Start Game!</PositiveButton>
        )}

        {isHost && <NegativeButton onClick={() => endGame()}>Cancel Game</NegativeButton>}
        {!isHost && <NegativeButton onClick={() => leaveGame()}>Leave Game</NegativeButton>}

      </Frame>
    </Container>
  );
};

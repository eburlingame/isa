import { useState } from "react";
import styled from "styled-components";
import { colors } from "../constants";

import { useActions } from "../hooks";

const Container = styled.div`
  width: 100%;
  height: 100%;

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

  transition: 0.3s background-color;
  background-color: ${colors.accent};
  color: ${colors.white};

  :hover {
    background-color: ${colors.accentDark};
    transition: 0.3s background-color;
  }
`;

const Line = styled.hr`
  margin-top: 3em;
  margin-bottom: 2em;
  border: 0.5px solid ${colors.midGray};
`;

const ErrorText = styled.div`
  margin-top: 0.25em;
  color: ${colors.red};
  font-weight: 600;
`;

const CreateGameForm = () => {
  const [playerName, setPlayerName] = useState("");

  const { createGame } = useActions();

  return (
    <FormContainer>
      <h1>Create a Game</h1>
      <Input
        placeholder="Your name"
        value={playerName}
        onChange={(e) => setPlayerName(e.target.value)}
        onKeyPress={(e) => {
          if (e.key === "Enter") {
            createGame(playerName);
          }
        }}
      />
      <Button onClick={() => createGame(playerName)}>Create Game</Button>
    </FormContainer>
  );
};

const JoinGameForm = () => {
  const [playerName, setPlayerName] = useState("");
  const [gameCode, setGameCode] = useState("");
  const [error, setError] = useState(null);

  const { joinGame } = useActions();

  const doJoin = async () => {
    if (playerName.length === 0) {
      setError("Invalid player name");
      return;
    }

    if (gameCode.length !== 4) {
      setError("Invalid game code");
      return;
    }

    try {
      await joinGame(playerName, gameCode);
    } catch (e) {
      setError(e);
    }
  };

  const onKeyPress = (e) => {
    if (e.key === "Enter") {
      doJoin();
    }
  };

  return (
    <FormContainer>
      <h1>Join a Game</h1>
      <Input
        placeholder="Your name"
        value={playerName}
        onChange={(e) => setPlayerName(e.target.value)}
        onKeyPress={onKeyPress}
      />
      <Input
        placeholder="Game code"
        value={gameCode}
        onChange={(e) => setGameCode(e.target.value)}
        onKeyPress={onKeyPress}
      />
      <Button onClick={doJoin}>Join Game</Button>
      <ErrorText>{error}</ErrorText>
    </FormContainer>
  );
};

export default ({}) => {
  return (
    <Container>
      <Frame>
        <CreateGameForm />
        <Line />
        <JoinGameForm />
      </Frame>
    </Container>
  );
};

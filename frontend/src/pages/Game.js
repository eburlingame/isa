import styled from "styled-components";

import Deck from "../components/Deck";
import Hand from "../components/Hand";
import Gameboard from "../components/Gameboard";
import { useEffect, useRef, useState } from "react";
import { isValidPlay, nextValidCard } from "../gameUtils";
import { useActions } from "../hooks";
import ColorPicker from "../components/ColorPicker";

import { colors } from "../constants";

const Container = styled.div`
  height: 100%;
  width: 100%;
  background-color: #fcfcfc;

  position: relative;
`;

const ButtonContainer = styled.div`
  position: absolute;
  right: 5%;
  top: 2%;
`;

const Button = styled.div`
  padding: 0.5em;

  background-color: ${colors.accent};
  transition: 0.3s background-color;

  color: ${colors.white};
  border-radius: 0.25em;

  cursor: pointer;

  :hover {
    background-color: ${colors.accentDark};
    transition: 0.3s background-color;
  }
`;

const DeckContainer = styled.div`
  position: absolute;
  left: 2%;
  width: 10%;
  minWidth: 200px;
`;

const ColorPickerContainer = styled.div`
  position: absolute;
  height: 200px;
  width: 50vw;
  max-width: 500px;
  bottom: 25%;
  left: 50%;
  transform: translate(-50%, 0%);
`;

const HandContainer = styled.div`
  position: absolute;
  height: 200px;
  width: 85vw;
  max-width: 1200px;
  bottom: 0%;
  left: 50%;
  transform: translate(-50%, 0%);
`;

const GameboardContainer = styled.div`
  position: absolute;
  left: 50%;
  top: 0%;
  transform: translate(-50%, 0%);
  height: calc(100vh - 200px);
  width: 85vh;
`;

const useCardSelection = (you, yourTurn, wildColor, discardPileTop, mustDraw) => {
  const [selectedCard, setSelectedCard] = useState(0);

  const trySelectingCard = (index) => {
    if (index === null) {
      return;
    }

    if (index === -1 || mustDraw > 0) {
      setSelectedCard(-1);
      return;
    }

    if (index >= you.cards.length || index < 0) {
      return;
    }

    if (isValidPlay(discardPileTop, wildColor, you.cards[index])) {
      setSelectedCard(index);
    }
  };

  useEffect(() => {
    const firstValidCard = nextValidCard(
      discardPileTop,
      wildColor,
      you.cards,
      -1,
      1
    );

    if (firstValidCard === null || mustDraw > 0) {
      setSelectedCard(-1);
    } else {
      setSelectedCard(firstValidCard);
    }
  }, [discardPileTop, you.cards]);

  return [selectedCard, trySelectingCard];
};

const COLORS = ["Y", "R", "G", "B"];

const nextColor = (color, dir) => {
  const index = COLORS.indexOf(color);
  const newIndex = index + dir;

  if (newIndex >= 0 && newIndex < COLORS.length) {
    console.log(COLORS[newIndex]);
    return COLORS[newIndex];
  }

  return color;
};

const useColorPicker = () => {
  const [pickingColor, setPickingColor] = useState(false);
  const [selectedColor, setSelectedColor] = useState("Y");

  const selectColor = (color) => {
    if (COLORS.includes(color)) {
      setSelectedColor(color);
    }
  };

  const pickColor = (value) => {
    setPickingColor(value);
  };

  return [pickingColor, selectedColor, pickColor, selectColor];
};

export default ({
  gameCode,
  game: {
    activePlayer,
    direction,
    discardPileTop,
    drawPileCount,
    mustDraw,
    otherPlayers,
    you,
    wildColor,
  },
  isHost,
}) => {
  const containerRef = useRef();

  const { playCard, drawCard, endGame, leaveGame } = useActions();

  const yourTurn = otherPlayers[activePlayer].name === you.name;

  const [
    pickingColor,
    selectedColor,
    pickColor,
    selectColor,
  ] = useColorPicker();

  const [selectedCard, trySelectingCard] = useCardSelection(
    you,
    yourTurn,
    wildColor,
    discardPileTop,
    mustDraw
  );

  const noValidCards =
    nextValidCard(discardPileTop, wildColor, you.cards, -1, 1) === null;

  const tryDrawCard = () => {
    if (!yourTurn) return;

    drawCard();
  };

  const tryPlayCard = (index) => {
    if (!yourTurn || index != selectedCard) return;

    if (pickingColor) {
      playCard(selectedCard, selectedColor);
      pickColor(false);
    } else if (selectedCard < 0) {
      drawCard();
    } else {
      if (you.cards[selectedCard].includes("wild")) {
        pickColor(true);
      } else {
        playCard(selectedCard, "");
      }
    }
  };

  // TODO: How to use React for this??
  const handleKeyPress = (e) => {
    if (!yourTurn) {
      return;
    }

    if (e.key === "Escape") {
      pickColor(false);
    }

    if (e.key === "ArrowRight") {
      if (pickingColor) {
        selectColor(nextColor(selectedColor, 1));
      } else {
        trySelectingCard(
          nextValidCard(discardPileTop, wildColor, you.cards, selectedCard, 1)
        );
      }
    }

    if (e.key === "ArrowLeft") {
      if (pickingColor) {
        selectColor(nextColor(selectedColor, -1));
      } else {
        trySelectingCard(
          nextValidCard(discardPileTop, wildColor, you.cards, selectedCard, -1)
        );
      }
    }

    if (e.key === " ") {
      tryPlayCard(selectedCard);
    }
  };
  document.body.onkeydown = handleKeyPress;

  return (
    <Container onKeyPress={handleKeyPress} ref={containerRef}>
      <ButtonContainer>
        {isHost && <Button onClick={() => endGame()}>End Game</Button>}
        {!isHost && <Button onClick={() => leaveGame()}>Leave Game</Button>}
      </ButtonContainer>

      <DeckContainer>
        <Deck
          cardCount={drawPileCount}
          selected={selectedCard === -1 && yourTurn}
          mustDraw={yourTurn ? mustDraw : 0}
          noValidCards={noValidCards}
          trySelectingCard={trySelectingCard}
          tryDrawCard={tryDrawCard}
        />
      </DeckContainer>

      <GameboardContainer>
        <Gameboard
          activePlayer={activePlayer}
          direction={direction}
          players={otherPlayers}
          wildColor={wildColor}
          discardPile={[discardPileTop]}
        />
      </GameboardContainer>

      <ColorPickerContainer>
        {pickingColor && (
          <ColorPicker
            selectedColor={selectedColor}
            selectColor={selectColor}
            playCard={() => tryPlayCard(selectedCard)}
          />
        )}
      </ColorPickerContainer>

      <HandContainer>
        <Hand
          yourTurn={yourTurn}
          cards={you.cards}
          selectedCard={selectedCard}
          trySelectingCard={trySelectingCard}
          playCard={tryPlayCard}
        />
      </HandContainer>
    </Container>
  );
};

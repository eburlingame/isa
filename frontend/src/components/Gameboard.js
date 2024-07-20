import styled from "styled-components";
import DiscardPile from "./DiscardPile";
import PlayerRing from "./PlayerRing";

const Container = styled.div`
  position: relative;
  height: 100%;
  width: 100%;
`;

export default ({
  activePlayer,
  direction,
  players,
  discardPile,
  wildColor,
}) => {
  return (
    <Container>
      <DiscardPile wildColor={wildColor} discardPile={discardPile} />

      <PlayerRing
        players={players}
        activePlayer={activePlayer}
        direction={direction}
      />
    </Container>
  );
};

import styled from "styled-components";
import { colors } from "../constants";
import { BlankCard } from "./Card";

const Container = styled.div`
  width: 100%;

  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
`;

const CardCount = styled.div`
  margin-top: 20px;

  font-weight: 600;
  font-size: 2em;
  margin-bottom: 0.75em;
  color: ${colors.darkGray};
`;

const CardLabel = styled.div`
  margin-top: 20px;

  font-weight: 600;
  font-size: 2em;
  margin-bottom: 0.75em;

  color: ${props => props.flashing ? colors.accent : colors.darkDarkGray};

  animation: ${props => props.flashing ? "2.4s flashing infinite" : "none"};
  @keyframes flashing {
    0% {
      color: ${colors.accent};
    };

    50% {
      color: ${colors.white};
    }
  }
`;

const CardsContainer = styled.div`
  position: relative;
  cursor: pointer;
  width: 100%;
  height: 100%;
`;

const TopCardContainer = styled.div`
  position: absolute;
  transition: 0.4s left;
  top: 0%;
  left: ${(props) => (props.selected ? "45%" : "0%")};
  width: 100%;
`;

export default ({
  cardCount,
  selected,
  mustDraw,
  trySelectingCard,
  tryDrawCard,
  noValidCards
}) => {
  return (
    <Container selectable={selected}>
      <CardCount>{cardCount} cards</CardCount>

      <CardsContainer
        onClick={tryDrawCard}
        onMouseOver={() => trySelectingCard(-1)}
      >
        <BlankCard selectable={true} />

        <TopCardContainer selected={selected}>
          <BlankCard selectable={true} />
        </TopCardContainer>
      </CardsContainer>

      {mustDraw === 0 && noValidCards && <CardLabel>Draw a card</CardLabel>}
      {mustDraw > 0 && <CardLabel flashing={true}>Draw {mustDraw} {mustDraw > 1 ? "cards" : "card"}</CardLabel>}
    </Container>
  );
}; 

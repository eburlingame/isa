import styled from "styled-components";
import { colors } from "../constants";
import { Card } from "./Card";

const Container = styled.div`
  width: 100%;
  height: 100%;

  background-color: ${colors.darkDarkGray};
  border-radius: 1em 1em 0px 0px;
  position: relative;
`;

const Curtain = styled.div`
  position: absolute;
  left: 0%;
  top: 0%;
  right: 0%;
  bottom: 0%;

  border-radius: 1em 1em 0px 0px;

  background-color: ${colors.black};
  opacity: 50%;
  z-index: 1000;
`;

const CardsContainer = styled.div`
  box-sizing: border-box;
  height: 100%;
`;

const cardWidth = "120px";
const edgeMargin = "15px";

const CardContainer = styled.div`
  margin: ${edgeMargin};

  cursor: pointer;
  box-sizing: border-box;
  position: absolute;

  left: calc(
    ${(props) => props.index / (props.count - 1)} *
      (100% - ${cardWidth} - 2 * ${edgeMargin})
  );
  z-index: ${(props) => props.index};

  top: ${(props) => (props.selected ? "-35%" : "0%")};
  transition: 0.4s top;

  width: ${cardWidth};
`;

export default ({
  yourTurn,
  selectedCard,
  cards,
  trySelectingCard,
  playCard,
}) => {
  return (
    <Container>
      <CardsContainer numCards={cards.length}>
        {cards.map((value, index) => (
          <CardContainer
            key={index}
            index={index}
            count={cards.length}
            selected={yourTurn && selectedCard === index}
            onMouseOver={() => trySelectingCard(index)}
            onClick={() => playCard(index)}
          >
            <Card value={value} />
          </CardContainer>
        ))}
      </CardsContainer>

      {!yourTurn && <Curtain />}
    </Container>
  );
};

import styled from "styled-components";
import { colors } from "../constants";
import { extractCardColor } from "../gameUtils";
import { Card } from "./Card";
import Player from "./Player";

const Container = styled.div`
  position: absolute;
  transform: translate(-50%, -50%);
  left: 50%;
  top: 50%;

  border-radius: 1000px;
  background-color: ${(props) => props.color};
  border: 1em solid ${(props) => props.darkColor};

  width: 30vh;
  height: 30vh;
`;

const PileContainer = styled.div`
  position: absolute;
  transform: translate(-50%, -50%);

  top: 50%;
  left: 50%;
  width: 50%;
  height: 70%;
`;

export default ({ discardPile, wildColor }) => {
  const { color, dark } = extractCardColor(discardPile[0], wildColor);

  return (
    <Container color={color} darkColor={dark}>
      {discardPile.length === 0 && <div>No cards</div>}

      {discardPile.length > 0 && (
        <PileContainer>
          <Card value={discardPile[0]} />
        </PileContainer>
      )}
    </Container>
  );
};

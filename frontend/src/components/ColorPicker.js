import styled from "styled-components";
import { colors } from "../constants";
import { BlankCard, Card } from "./Card";

const Container = styled.div`
  width: 100%;
  height: 100%;

  background-color: ${colors.darkDarkGray};
  border-radius: 1em;
  position: relative;

  display: flex;
  align-items: center;
  justify-content: space-between;

  padding: 2em;
`;

const CardContainer = styled.div`
  height: 100%;

  cursor: pointer;

  transition: 0.4s width, 0.4s transform;

  width: ${(props) => (props.selected ? "25%" : "20%")};
  transform: ${(props) =>
    props.selected ? "translateY(-20%)" : "translateY(0%)"};
`;

const Text = styled.div`
    color; ${(props) => colors.white}
`;

const colorsValues = ["Y", "R", "G", "B"];

export default ({ selectedColor, selectColor, playCard }) => {
  return (
    <Container>
      {colorsValues.map((color) => (
        <CardContainer
          key={color}
          selected={selectedColor === color}
          onMouseOver={() => selectColor(color)}
          onClick={playCard}
        >
          <Card value={color} />
        </CardContainer>
      ))}
    </Container>
  );
};

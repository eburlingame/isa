import styled from "styled-components";
import { colors } from "../constants";
import { range } from "lodash";
import CardSet from "./CardSet";

const Container = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 15vh;
  padding-left: 5%;
  padding-right: 5%;

  border: ${(props) =>
    props.playerActive || props.animate
      ? "5px solid " + (props.animate ? colors.green : colors.accent)
      : "0px"};

  animation: ${(props) =>
    props.animate ? "highlight 2.5s linear infinite" : "none"};

  @keyframes highlight {
    0% {
      border: ${(props) => "5px solid " + colors.green};
    }

    50% {
      border: ${(props) => "5px solid " + colors.lightGray};
    }
  }

  border-radius: 1em;

  background-color: ${(props) =>
    props.playerActive ? colors.lightGray : colors.white};
`;

const Name = styled.div`
  font-size: 2em;
  margin-bottom: 0.25em;
  font-weight: 900;
  color: ${(props) => (props.playerActive ? colors.accent : colors.darkGray)};
`;

export default ({ name, numCards, playerActive }) => {
  return (
    <Container playerActive={playerActive} animate={numCards === 1}>
      <Name playerActive={playerActive}>{name}</Name>

      <CardSet numCards={numCards} />
    </Container>
  );
};

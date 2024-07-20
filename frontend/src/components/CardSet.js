import { BlankCard } from "./Card";
import styled from "styled-components";
import { range } from "lodash";

const CardSetContainer = styled.div`
  width: 100%;
  padding-top: 80%;
  position: relative;
`;

const CardContainer = styled.div`
  position: absolute;
  width: 50%;
  left: ${(props) => (props.index / (props.numCards - 1)) * 50}%;
  top: 0%;
  z-index: ${(props) => props.zIndex};
`;

export default ({ numCards }) => {
  return (
    <CardSetContainer>
      {range(numCards).map((_, index) => (
        <CardContainer
          key={index}
          index={index}
          numCards={numCards}
          zIndex={index}
        >
          <BlankCard />
        </CardContainer>
      ))}
    </CardSetContainer>
  );
};

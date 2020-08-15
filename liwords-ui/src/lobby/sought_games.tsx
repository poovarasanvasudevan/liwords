/* eslint-disable jsx-a11y/click-events-have-key-events */
/* eslint-disable jsx-a11y/no-noninteractive-element-interactions */
import { Table } from 'antd';
import React, { ReactNode } from 'react';
import { challRuleToStr, timeCtrlToDisplayName } from '../store/constants';
import { SoughtGame } from '../store/reducers/lobby_reducer';
import { FundOutlined } from '@ant-design/icons/lib';

type SoughtGameProps = {
  game: SoughtGame;
  newGame: (seekID: string) => void;
  userID?: string;
  username?: string;
  requests: SoughtGame[];
};

/* const SoughtGameItem = (props: SoughtGameProps) => {
  const { game, userID, username } = props;
  const [tt, tc] = timeCtrlToDisplayName(
    game.initialTimeSecs,
    game.incrementSecs,
    game.maxOvertimeMinutes
  );
  let innerel = (
    // eslint-disable-next-line jsx-a11y/no-static-element-interactions
    <div
      className="game"
      onClick={(event: React.MouseEvent) => props.newGame(game.seekID)}
    >
      <RatingBadge rating={game.userRating} player={game.seeker} />
      {game.lexicon} ({game.rated ? 'Rated' : 'Casual'})(
      {`${game.initialTimeSecs / 60} / ${game.incrementSecs}`})
      <Tag color={tc}>{tt}</Tag>
      {challRuleToStr(game.challengeRule)}
      {game.incrementSecs === 0
        ? ` (Max OT: ${game.maxOvertimeMinutes} min.)`
        : ''}
    </div>
  );

  if (game.receiver && game.receiver.getUserId() !== '') {
    console.log('reciever', game.receiver);

    if (userID === game.receiver.getUserId()) {
      // This is the receiver of the match request.
      innerel = (
        <Badge.Ribbon text="Match Request" color="volcano">
          {innerel}
        </Badge.Ribbon>
      );
    } else {
      // We must be the sender of the match request.
      if (username !== game.seeker) {
        throw new Error(`unexpected seeker${username}, ${game.seeker}`);
      }
      innerel = (
        <Badge.Ribbon
          text={`Outgoing Match Request to ${game.receiver.getDisplayName()}`}
        >
          {innerel}
        </Badge.Ribbon>
      );
    }
  }

  return <li>{innerel}</li>;
};
*/

type Props = {
  isMatch?: boolean;
  newGame: (seekID: string) => void;
  userID?: string;
  username?: string;
  requests: Array<SoughtGame>;
};

export const SoughtGames = (props: Props) => {
  // @ts-ignore
  // @ts-ignore
  const columns = [
    {
      title: 'Player',
      className: 'seeker',
      dataIndex: 'seeker',
      key: 'seeker',
    },
    {
      title: 'Rating',
      className: 'rating',
      dataIndex: 'rating',
      key: 'rating',
      sorter: (a: SoughtGameTableData, b: SoughtGameTableData) =>
        parseInt(a.rating) - parseInt(b.rating),
    },
    {
      title: 'Words',
      className: 'lexicon',
      dataIndex: 'lexicon',
      key: 'lexicon',
      filters: [
        {
          text: 'CSW19',
          value: 'CSW19',
        },
        {
          text: 'NWL18',
          value: 'NWL18',
        },
      ],
      filterMultiple: false,
      onFilter: (
        value: string | number | boolean,
        record: SoughtGameTableData
      ) => record.lexicon.indexOf(value.toString()) === 0,
    },
    {
      title: 'Time',
      className: 'time',
      dataIndex: 'time',
      key: 'time',

      sorter: (a: SoughtGameTableData, b: SoughtGameTableData) =>
        a.time.localeCompare(b.time),
    },
    {
      title: 'Details',
      className: 'details',
      dataIndex: 'details',
      key: 'details',
    },
  ];

  type SoughtGameTableData = {
    seeker: string;
    rating: string;
    lexicon: string;
    time: string;
    details?: ReactNode;
    seekID: string;
  };

  const formatGameData = (games: SoughtGame[]): SoughtGameTableData[] => {
    const gameData: SoughtGameTableData[] = games.map(
      (sg: SoughtGame): SoughtGameTableData => {
        const challengeFormat = (cr: number) => {
          return (
            <span className={`challenge-rule mode_${challRuleToStr(cr)}`}>
              {challRuleToStr(cr)}
            </span>
          );
        };
        const timeFormat = (
          initialTimeSecs: number,
          incrementSecs: number,
          maxOvertime: number
        ): string => {
          // TODO: Pull in from time control in seek window
          const label = timeCtrlToDisplayName(
            initialTimeSecs,
            incrementSecs,
            maxOvertime
          )[0];
          return `${label} ${initialTimeSecs / 60}/${incrementSecs}`;
        };
        const getDetails = () => {
          return (
            <>
              {challengeFormat(sg.challengeRule)}
              {sg.rated ? <FundOutlined title="Rated" /> : null}
            </>
          );
        };
        return {
          seeker: sg.seeker,
          rating: sg.userRating,
          lexicon: sg.lexicon,
          time: timeFormat(
            sg.initialTimeSecs,
            sg.incrementSecs,
            sg.maxOvertimeMinutes
          ),
          details: getDetails(),
          seekID: sg.seekID,
        };
      }
    );
    return gameData;
  };

  return (
    <>
      {props.isMatch ? <h4>Incoming Requests</h4> : <h4>Available Games</h4>}
      <Table
        className={`games ${props.isMatch ? 'match' : 'seek'}`}
        dataSource={formatGameData(props.requests)}
        columns={columns}
        pagination={{
          hideOnSinglePage: true,
        }}
        onRow={(record, rowIndex) => {
          return {
            onClick: (event) => {
              props.newGame(record.seekID);
            },
          };
        }}
        rowClassName="game-listing"
      />
    </>
  );
};

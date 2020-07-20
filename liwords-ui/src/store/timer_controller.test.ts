import { millisToTimeStr } from './timer_controller';

it('tests millis to time', () => {
  expect(millisToTimeStr(479900)).toEqual('08:00');
  expect(millisToTimeStr(479000)).toEqual('07:59');
  expect(millisToTimeStr(60000)).toEqual('01:00');
  expect(millisToTimeStr(59580)).toEqual('01:00');
  expect(millisToTimeStr(8900)).toEqual('00:08.9');
  expect(millisToTimeStr(890)).toEqual('00:00.9');
  expect(millisToTimeStr(-1)).toEqual('-00:00.0');
  expect(millisToTimeStr(-40)).toEqual('-00:00.0');
  expect(millisToTimeStr(-89)).toEqual('-00:00.1');
  expect(millisToTimeStr(-890)).toEqual('-00:00.9');
  expect(millisToTimeStr(-990)).toEqual('-00:01.0');
  expect(millisToTimeStr(-8900)).toEqual('-00:08');
  expect(millisToTimeStr(-10000)).toEqual('-00:10');
  expect(millisToTimeStr(-10300)).toEqual('-00:10');
  expect(millisToTimeStr(-10600)).toEqual('-00:10');
  expect(millisToTimeStr(-11000)).toEqual('-00:11');
  expect(millisToTimeStr(-59000)).toEqual('-00:59');
  expect(millisToTimeStr(-59700)).toEqual('-00:59');
  expect(millisToTimeStr(-60000)).toEqual('-01:00');
});

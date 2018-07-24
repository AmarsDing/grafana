import * as MultiStatPanel from './types';

const panelDefaults: MultiStatPanel.PanelOptions = {
  links: [],
  datasource: null,
  maxDataPoints: 100,
  interval: null,
  targets: [{}],
  cacheTimeout: null,
  format: 'none',
  valueName: 'avg',
  prefixFontSize: '50%',
  valueFontSize: '80%',
  postfixFontSize: '50%',
  thresholds: [],
  colorBackground: false,
  colorValue: false,
  sparkline: {
    show: false,
    // full: false,
    // lineColor: 'rgb(31, 120, 193)',
    // fillColor: 'rgba(31, 118, 189, 0.18)',
  },
  layout: 'horizontal',
  viewMode: 'stats',
};

const valueNameOptions = [
  { value: 'min', text: 'Min' },
  { value: 'max', text: 'Max' },
  { value: 'avg', text: 'Average' },
  { value: 'current', text: 'Current' },
  { value: 'total', text: 'Total' },
  { value: 'name', text: 'Name' },
  { value: 'first', text: 'First' },
  { value: 'delta', text: 'Delta' },
  { value: 'diff', text: 'Difference' },
  { value: 'range', text: 'Range' },
  { value: 'last_time', text: 'Time of last point' },
];

const layoutOptions = [{ value: 'horizontal', text: 'Horizontal' }, { value: 'vertical', text: 'Vertical' }];

const viewModeOptions = [{ value: 'stats', text: 'Stats' }, { value: 'bars', text: 'Bars' }];

const defaults = {
  panelDefaults,
  valueNameOptions,
  layoutOptions,
  viewModeOptions,
};

export default defaults;

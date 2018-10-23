import expect from 'expect';
import Enzyme from 'enzyme';
import Adapter from 'enzyme-adapter-react-15';
import {extractNonTranslatableParts, fillTranslatedText} from './actions';

Enzyme.configure({ adapter: new Adapter() });

test('Extract Non translatable Parts', () => {
    const output = '';
    expect(output.text).toEqual('');
});

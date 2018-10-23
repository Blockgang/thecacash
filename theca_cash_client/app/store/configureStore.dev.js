import { applyMiddleware, createStore} from 'redux';
import thecaClientApp from '../reducers/reducers';
import thunkMiddleware from 'redux-thunk';
import api from '../middleware/api'

const createStoreWithMiddleware = applyMiddleware(thunkMiddleware, api)(createStore);

const store = createStoreWithMiddleware(thecaClientApp);

export function configureStore() {
    return store;
}

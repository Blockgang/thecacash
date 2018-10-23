import { applyMiddleware, createStore} from 'redux';
import challengesApp from '../reducers/reducers';
import thunkMiddleware from 'redux-thunk';
import api from '../middleware/api'

const createStoreWithMiddleware = applyMiddleware(thunkMiddleware, api)(createStore);

const store = createStoreWithMiddleware(challengesApp);

export function configureStore() {
    return store;
}

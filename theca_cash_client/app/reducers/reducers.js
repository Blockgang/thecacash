import { combineReducers } from 'redux';
import {
    LOGIN_REQUEST, LOGIN_SUCCESS, LOGIN_FAILURE, LOGOUT_SUCCESS,
    ERROR_OCCURRED
} from '../actions/actions';

// The auth reducer. The starting state sets authentication
// based on a token being in local storage. In a real app,
// we would also want a util to check if the token is expired.

function auth(state = {
    isFetching: false,
    isAuthenticated: false
}, action) {
    switch (action.type) {
        case LOGIN_REQUEST:
            return Object.assign({}, state, {
                isFetching: true,
                isAuthenticated: false,
                user: action.credentials
            });
        case LOGIN_SUCCESS:
            return Object.assign({}, state, {
                isFetching: false,
                isAuthenticated: action.isAuthenticated,
                errorMessage: '',
            });
        case LOGIN_FAILURE:
            return Object.assign({}, state, {
                isFetching: false,
                isAuthenticated: false,
                errorMessage: action.message
            });
        case LOGOUT_SUCCESS:
            return Object.assign({}, state, {
                isFetching: true,
                isAuthenticated: false
            });
        default:
            return state;
    }
}

function messageHandler(state = {
    hasError: false,
    errorContext: '',
    message: ''
}, action) {
    switch (action.type) {
        case ERROR_OCCURRED:
            return Object.assign({}, state, {
                message: action.message === undefined ? state.message : action.message,
                errorContext: action.errorContext === undefined ? state.errorContext : action.errorContext,
                hasError: action.hasError === undefined ? state.hasError : action.hasError,
            });
        default:
            return state;
    }
}

const thecaClientApp = combineReducers({
    auth,
    messageHandler
});

export default thecaClientApp;

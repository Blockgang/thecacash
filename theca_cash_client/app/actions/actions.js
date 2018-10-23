// eslint-disable-next-line no-unused-vars
import {CALL_API} from '../middleware/api';

const SERVER_API_ADDRESS = process.env.SERVER_API_ADDRESS || 'localhost';
const SERVER_API_PORT = process.env.SERVER_API_PORT || 8000;
const SERVER_API_PROTOCOL = process.env.SERVER_API_PROTOCOL || 'http';

const SERVER_URL = SERVER_API_PROTOCOL + '://' + SERVER_API_ADDRESS + ':' + SERVER_API_PORT;
// There are three possible states for our login
// process and we need actions for each of them
export const LOGIN_REQUEST = 'LOGIN_REQUEST';
export const LOGIN_SUCCESS = 'LOGIN_SUCCESS';
export const LOGIN_FAILURE = 'LOGIN_FAILURE';

function requestLogin(creds) {
    return {
        type: LOGIN_REQUEST,
        isFetching: true,
        isAuthenticated: false,
        creds
    };
}

function receiveLogin(user) {
    return {
        type: LOGIN_SUCCESS,
        isFetching: false,
        isAuthenticated: true,
        id_token: user.id_token
    };
}

function loginError(message) {
    return {
        type: LOGIN_FAILURE,
        isFetching: false,
        isAuthenticated: false,
        message
    };
}

// Three possible states for our logout process as well.
// Since we are using JWTs, we just need to remove the token
// from localStorage. These actions are more useful if we
// were calling the API to log the user out
export const LOGOUT_REQUEST = 'LOGOUT_REQUEST';
export const LOGOUT_SUCCESS = 'LOGOUT_SUCCESS';
export const LOGOUT_FAILURE = 'LOGOUT_FAILURE';

function requestLogout() {
    return {
        type: LOGOUT_REQUEST,
        isFetching: true,
        isAuthenticated: true
    };
}

function receiveLogout() {
    return {
        type: LOGOUT_SUCCESS,
        isFetching: false,
        isAuthenticated: false
    };
}

// Calls the API to get a token and
// dispatches actions along the way
export function loginUser(creds) {
    const config = {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        data: `Username=${creds.username}&PasswordHash=${creds.password}`
    };

    return dispatch => {
        // We dispatch requestLogin to kickoff the call to the API
        dispatch(requestLogin(creds));
        return fetch(SERVER_URL + '/api/login', config)
            .then(response =>
                response.json()
                    .then(user => ({user, response}))
            ).then(({user, response}) => {
                if (!response.ok) {
                    // If there was a problem, we want to
                    // dispatch the error condition
                    dispatch(loginError(user.message));
                    // eslint-disable-next-line no-undef
                    return Promise.reject(user);
                }
                localStorage.setItem('authenticated', user.signup);
                dispatch(receiveLogin(user));
                return null;
            }).catch(err => console.log('Error: ', err));
    };
}

// Logs the user out
export function logoutUser() {
    return dispatch => {
        dispatch(requestLogout());
        localStorage.removeItem('authenticated');
        dispatch(receiveLogout());
    };
}

export function signupUser(creds) {
    const config = {
        method: 'POST',
        data: `Username=${creds.username}&PasswordHash=${creds.password}&EncryptedPK=${creds.pk}`
    };

    return dispatch => {
        // We dispatch requestLogin to kickoff the call to the API
        dispatch(requestLogin(creds));
        return fetch(SERVER_URL + '/api/signup', config)
            .then(response =>
                response.json()
                    .then(user => ({user, response}))
            ).then(({user, response}) => {
                if (!response.ok) {
                    // If there was a problem, we want to
                    // dispatch the error condition
                    dispatch(loginError(user.message));
                    // eslint-disable-next-line no-undef
                    return Promise.reject(user);
                }
                localStorage.setItem('authenticated', user.signup);
                dispatch(receiveLogin(user));
                return null;
            }).catch(err => console.log('Error: ', err));
    };
}

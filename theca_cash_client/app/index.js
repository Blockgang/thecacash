import React from 'react';
import { render } from 'react-dom';
import { AppContainer } from 'react-hot-loader';
import { configureStore, history } from './store/configureStore';
import Root from './containers/Root';
import 'semantic-ui/dist/semantic.min.css';

const store = configureStore();

const rootElement = document.getElementById('root');

render(
    <AppContainer>
        <Root store={store} history={history} />
    </AppContainer>,
    rootElement
);

if (module.hot) {
    module.hot.accept('./containers/Root', () => {
        const newConfigureStore = require('./store/configureStore');
        const newStore = newConfigureStore.configureStore();
        const newHistory = newConfigureStore.history;
        const NewRoot = require('./containers/Root').default;
        render(
            <AppContainer>
                <NewRoot store={newStore} history={newHistory} />
            </AppContainer>,
            rootElement
        );
    });
}

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import {Input, FormField, } from 'semantic-ui-react';

export default class SimpleInputField extends Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
        const {focused} = this.props;
        if(focused) {
            this.textInput.focus();
        }
    }

    render() {
        const { errorMessage, onChange, className, size, type, id, labelText, placeholder, required, onBlur, value, ...props} = this.props;
        return (
            <FormField>
                <Input
                    onBlur={onBlur}
                    onChange={onChange}
                    fluid
                    value={value}
                    size={size}
                    className={className}
                    id={id}
                    label={labelText}
                    ref={(input) => {this.textInput = input; }}
                    placeholder={placeholder}
                    name={name}
                    required={required}
                    type={type}
                    error={errorMessage}
                    {...props}
                />
            </FormField>
        );
    }
}

SimpleInputField.propTypes = {
    className: PropTypes.string,
    onBlur: PropTypes.func,
    onChange: PropTypes.func.isRequired,
    id: PropTypes.string.isRequired,
    labelText: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    reference: PropTypes.string,
    placeholder: PropTypes.string.isRequired,
    required: PropTypes.bool,
    focused: PropTypes.bool,
    type: PropTypes.string,
    size: PropTypes.string,
    value: PropTypes.string,
    errorMessage: PropTypes.string
};

SimpleInputField.defaultProps = {
    className: 'simpleInputField',
    required: true,
    type: 'text',
    size: 'huge',
    value: '',
    focused: false
};

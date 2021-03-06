import React, { Component, PropTypes } from 'react'
import { Link } from 'react-router'
import { I18n, Translate } from 'react-redux-i18n'
import { Field, propTypes as formPropTypes } from 'redux-form'
import cx from 'classnames'

import { createEmptyPromise } from 'utils'
import { Button, FieldText } from 'uis'
import { ACCOUNT_BASE_PATH, SIGN_IN_PATH } from '../index'

import './sign-up.view.styl'

export class SignUp extends Component {

  static propTypes = {
    className: PropTypes.string,
    location: PropTypes.object,
    actions: PropTypes.object,
    ...formPropTypes
  }

  getRootClassNames () {
    return cx(
      'signUpView',
      this.props.className
    )
  }

  handleSubmit = (values) => {
    const {
      location,
      actions
    } = this.props

    const formPromise = createEmptyPromise()

    actions.signUpUser({
      username: values.username,
      password: values.password,
      location,
      formPromise
    })

    return formPromise
  }

  renderSignUpForm () {
    const { handleSubmit, anyTouched, valid, submitting } = this.props

    return (
      <form onSubmit={handleSubmit(this.handleSubmit)}>
        <Field
          name='username'
          component={FieldText}
          placeholder={I18n.t('account.username')}
          autoFocus
        />
        <Field
          name='password'
          component={FieldText}
          type={'password'}
          placeholder={I18n.t('account.password')}
        />
        <Field
          name='passwordRetype'
          component={FieldText}
          type={'password'}
          placeholder={I18n.t('account.passwordRetype')}
        />
        <Button
          className={'SignUpHandler'}
          type={'primary'}
          htmlType={'submit'}
          disabled={!anyTouched || !valid}
          loading={submitting}
          block
        >
          <Translate value={'account.signUp'} />
        </Button>
      </form>
    )
  }

  render () {
    const { location } = this.props

    return (
      <div className={this.getRootClassNames()}>
        <div className={'accountLayoutForm'}>
          {this.renderSignUpForm()}
        </div>
        <Link
          className={'accountLayoutFooter'}
          to={`${ACCOUNT_BASE_PATH}/${SIGN_IN_PATH}${location.search}`}
        >
          <Translate value={'account.signInTip'} />
        </Link>
      </div>
    )
  }

}

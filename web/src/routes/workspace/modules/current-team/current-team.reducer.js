import { createAction, handleActions } from 'redux-actions'

const initialState = {
  teamId: null
}

export const setCurrentTeamAction = createAction('SET_CURRENT_TEAM')

export const currentTeamReducer = handleActions({

  [`${setCurrentTeamAction}`]: (state, action) => {
    if (!action.payload) {
      return state
    }

    const { teamId } = action.payload

    if (typeof teamId === undefined) {
      return state
    }

    return {
      ...state,
      teamId
    }
  }

}, initialState)
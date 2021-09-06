Rails.application.routes.draw do
  get "auth/callback"
  get "auth/failure"
  get "auth/logout"

  get "session/new"
  post "session/create"
  delete "session/destroy"

  get "user/index"
  get "user/show"
  get "user/new"
  post "user/create"
  delete "user/destroy"

  get "user/detail/show"

  # For details on the DSL available within this file, see https://guides.rubyonrails.org/routing.html

  root "user#index"
end

Rails.application.routes.draw do
  post 'session/create'
  delete 'session/destroy'

  get "user/index"
  get "user/show"
  post "user/create"
  patch "user/update"
  delete "user/destroy"
  # For details on the DSL available within this file, see https://guides.rubyonrails.org/routing.html

  root "user#index"
end

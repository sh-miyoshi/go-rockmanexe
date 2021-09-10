class Session < ApplicationRecord
  validates :session_name, presence: true
  validates :owner_id:presence, true
  validates :guest_id, presence: true
end

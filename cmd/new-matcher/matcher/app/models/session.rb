class Session < ApplicationRecord
  validates :name, presence: true
  validates :owner_id, presence: true
  validates :guest_id, presence: true
end

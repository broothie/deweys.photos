require 'dotenv/load'
require 'sinatra'
require 'sinatra/flash'
require 'sinatra/reloader' if development?
require 'sassc'
require 'google/cloud/firestore'
require 'google/cloud/storage'

enable :sessions
set session_secret: ENV.fetch('SESSION_SECRET') { SecureRandom.hex(64) }

not_found do
  redirect '/'
end

get '/' do
  @entries = all_entries
  title 'Home'
  erb :index
end

get '/manage' do
  require_logged_in!

  @entries = all_entries
  title 'Manage Photos'
  erb :manage
end

post '/photos' do
  require_logged_in!

  # Get uploaded photo data
  photo = params[:photo]
  filename = photo[:filename]
  local_file = photo[:tempfile]

  # Upload file to cloud
  cloud_file = bucket.create_file(local_file, filename)

  # Save entry to db
  entry = entries.add(
    title: params[:title],
    blurb: params[:blurb],
    filename: filename,
    url: cloud_file.public_url,
  )

  # Show uploaded photo on homepage
  redirect "/#photos-#{entry.document_id}"
end

get '/photos/:entry_id' do |entry_id|
  require_logged_in!

  @entry = entries.doc(entry_id).get
  title 'Edit Photo'
  erb :edit_photo
end

post '/photos/:entry_id' do |entry_id|
  require_logged_in!

  entries.doc(entry_id).update({
    title: params[:title],
    blurb: params[:blurb]
  })

  redirect "/#photos-#{entry_id}"
end

get '/login' do
  redirect_to_callback if logged_in?

  title 'Log In'
  erb :login
end

post '/login' do
  if password_correct?(params[:password])
    log_in!
    redirect_to_callback
  else
    flash[:error] = 'Incorrect password'
    redirect '/login'
  end
end

delete '/login' do
  log_out!
end

get '/style.css' do
  scss :base
end

helpers do
  # Template
  def title(title)
    @title = title
  end

  # Login
  def password_correct?(password)
    (@password ||= ENV['PASSWORD']) == password
  end

  def log_in!
    session[:logged_in] = true
  end

  def log_out!
    session.delete(:logged_in)
  end

  def logged_in?
    session[:logged_in]
  end

  def require_logged_in!
    callback_redirect '/login' unless logged_in?
  end

  # Callbacks
  def callback_redirect(path)
    session[:callback] = request.path
    redirect path
  end

  def redirect_to_callback
    redirect session.delete(:callback) || '/'
  end

  # GCP
  def all_entries
    entries.get.sort_by(&:created_at).reverse
  end

  def entries
    @entries ||= Google::Cloud::Firestore
      .new(**gcp_args)
      .col(settings.development? ? 'entries-development' : 'entries')
  end

  def bucket
    @bucket ||= Google::Cloud::Storage
      .new(**gcp_args)
      .bucket(settings.development? ? 'deweys-photos-development' : 'deweys-photos')
  end

  def gcp_args
    @gcp_args ||= {
      project: 'deweys-photos',
      keyfile: settings.development? ? 'gcloud-key.json' : nil
    }
  end
end

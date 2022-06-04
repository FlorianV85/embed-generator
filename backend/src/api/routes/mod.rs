pub use auth_exchange::route_auth_exchange;
pub use auth_redirect::route_auth_redirect;
pub use guild_channel_list::route_guild_channel_list;
pub use guild_emoji_list::route_guild_emoji_list;
pub use guild_get::route_guild_get;
pub use guild_list::route_guild_list;
pub use guild_role_list::route_guild_role_list;
pub use guild_sticker_list::route_guild_sticker_list;
pub use links::route_link_discord;
pub use links::route_link_invite;
pub use links::route_link_source;
pub use message_create::route_message_create;
pub use message_delete::route_message_delete;
pub use message_update::route_message_update;
#[cfg(feature = "frontend")]
pub use serve_frontend::route_serve_frontend;
pub use user_get_me::route_user_get_me;

mod auth_exchange;
mod auth_redirect;
mod guild_channel_list;
mod guild_emoji_list;
mod guild_get;
mod guild_list;
mod guild_role_list;
mod guild_sticker_list;
mod links;
mod message_create;
mod message_delete;
mod message_get;
mod message_list;
mod message_send;
mod message_update;
mod user_get_me;
#[cfg(feature = "frontend")]
mod serve_frontend;

module Twirp.JSONPB
    -- * Typeclasses
  ( FromJSONPB(..)
  , ToJSONPB(..)
    -- * Operators
  , (.:)
  , (.=)
    -- * Options for controlling codec behavior (e.g., emitting default-valued
    --   fields in JSON payloads)
  , Options(..)
  , defaultOptions
    -- * JSONPB codec entry points
  , eitherDecode
  , encode
    -- * Helper functions
  -- , enumFieldEncoding
  -- , enumFieldString
  , object
  , pair
  , pairs
  , parseField
  , toAesonEncoding
  , toAesonValue
    -- * Aeson re-exports
  , A.Value(..)
  , A.ToJSON(..)
  , A.FromJSON(..)
  , A.typeMismatch
  , A.withObject
  ) where

import qualified Data.Aeson as A
import qualified Data.Aeson.Types as A
import           Twirp.JSONPB.Class

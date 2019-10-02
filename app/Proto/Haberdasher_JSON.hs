-- Code generated by protoc-gen-jsonpb_haskell 0.1.0, DO NOT EDIT.
{-# OPTIONS_GHC -Wno-orphans -Wno-unused-imports -Wno-missing-export-lists #-}
module Proto.Haberdasher_JSON where

import           Prelude(($), (.), (<$>), pure, show)

import           Data.ProtoLens.Runtime.Lens.Family2 ((^.), (.~), (&))
import           Data.Monoid (mconcat)
import           Control.Monad (msum)
import           Data.ProtoLens (defMessage)
import qualified Data.Aeson as A
import qualified Data.Aeson.Encoding as E
import           Twirp.JSONPB as JSONPB
import qualified Data.Text as T

import           Proto.Haberdasher as P
import           Proto.Haberdasher_Fields as P

instance FromJSONPB Size where
  parseJSONPB = withObject "Size" $ \obj -> do
    inches' <- obj .: "inches"
    pure $ defMessage
      & P.inches .~ inches'

instance ToJSONPB Size where
  toJSONPB x = object
    [ "inches" .= (x^.inches)
    ]
  toEncodingPB x = pairs
    [ "inches" .= (x^.inches)
    ]

instance FromJSON Size where
  parseJSON = parseJSONPB

instance ToJSON Size where
  toJSON = toAesonValue
  toEncoding = toAesonEncoding

instance FromJSONPB Hat where
  parseJSONPB = withObject "Hat" $ \obj -> do
    inches' <- obj .: "inches"
    color' <- obj .: "color"
    name' <- obj .: "name"
    pure $ defMessage
      & P.inches .~ inches'
      & P.color .~ color'
      & P.name .~ name'

instance ToJSONPB Hat where
  toJSONPB x = object
    [ "inches" .= (x^.inches)
    , "color" .= (x^.color)
    , "name" .= (x^.name)
    ]
  toEncodingPB x = pairs
    [ "inches" .= (x^.inches)
    , "color" .= (x^.color)
    , "name" .= (x^.name)
    ]

instance FromJSON Hat where
  parseJSON = parseJSONPB

instance ToJSON Hat where
  toJSON = toAesonValue
  toEncoding = toAesonEncoding

instance FromJSONPB Bill'Extra where
  parseJSONPB = A.withObject "Bill'Extra" $ \obj -> mconcat
    [
      Bill'VatInfo <$> parseField obj "vat_info"
    , Bill'ZipCode <$> parseField obj "zip_code"
    ]

instance ToJSONPB Bill'Extra where
  toJSONPB (Bill'VatInfo x) = object [ "vat_info" .= x ]
  toJSONPB (Bill'ZipCode x) = object [ "zip_code" .= x ]
  toEncodingPB (Bill'VatInfo x) = pairs [ "vat_info" .= x ]
  toEncodingPB (Bill'ZipCode x) = pairs [ "zip_code" .= x ]

instance FromJSON Bill'Extra where
  parseJSON = parseJSONPB

instance ToJSON Bill'Extra where
  toJSON = toAesonValue
  toEncoding = toAesonEncoding

instance FromJSONPB Bill where
  parseJSONPB = withObject "Bill" $ \obj -> do
    price' <- obj A..:? "price"
    status' <- obj .: "status"
    extra' <- obj A..:? "extra"
    pure $ defMessage
      & P.maybe'price .~ price'
      & P.status .~ status'
      & P.maybe'extra .~ extra'

instance ToJSONPB Bill where
  toJSONPB x = object
    [ "price" .= (x^.maybe'price)
    , "status" .= (x^.status)
    , "extra" .= (x^.maybe'extra)
    ]
  toEncodingPB x = pairs
    [ "price" .= (x^.maybe'price)
    , "status" .= (x^.status)
    , "extra" .= (x^.maybe'extra)
    ]

instance FromJSON Bill where
  parseJSON = parseJSONPB

instance ToJSON Bill where
  toJSON = toAesonValue
  toEncoding = toAesonEncoding

instance FromJSONPB Bill'BillingStatus where
  parseJSONPB (JSONPB.String "UN_PAID") = pure Bill'UN_PAID
  parseJSONPB (JSONPB.String "PAID") = pure Bill'PAID
  parseJSONPB x = typeMismatch "BillingStatus" x

instance ToJSONPB Bill'BillingStatus where
  toJSONPB x _ = A.String . T.toUpper . T.pack $ show x
  toEncodingPB x _ = E.text . T.toUpper . T.pack  $ show x

instance FromJSON Bill'BillingStatus where
  parseJSON = parseJSONPB

instance ToJSON Bill'BillingStatus where
  toJSON = toAesonValue
  toEncoding = toAesonEncoding

instance FromJSONPB Test where
  parseJSONPB = withObject "Test" $ \obj -> do
    items' <- obj .: "items"
    altPrices' <- obj .: "altPrices"
    pure $ defMessage
      & P.items .~ items'
      & P.altPrices .~ altPrices'

instance ToJSONPB Test where
  toJSONPB x = object
    [ "items" .= (x^.items)
    , "altPrices" .= (x^.altPrices)
    ]
  toEncodingPB x = pairs
    [ "items" .= (x^.items)
    , "altPrices" .= (x^.altPrices)
    ]

instance FromJSON Test where
  parseJSON = parseJSONPB

instance ToJSON Test where
  toJSON = toAesonValue
  toEncoding = toAesonEncoding

instance FromJSONPB Price where
  parseJSONPB = withObject "Price" $ \obj -> do
    dollars' <- obj .: "dollars"
    cents' <- obj .: "cents"
    pure $ defMessage
      & P.dollars .~ dollars'
      & P.cents .~ cents'

instance ToJSONPB Price where
  toJSONPB x = object
    [ "dollars" .= (x^.dollars)
    , "cents" .= (x^.cents)
    ]
  toEncodingPB x = pairs
    [ "dollars" .= (x^.dollars)
    , "cents" .= (x^.cents)
    ]

instance FromJSON Price where
  parseJSON = parseJSONPB

instance ToJSON Price where
  toJSON = toAesonValue
  toEncoding = toAesonEncoding

instance FromJSONPB Ping where
  parseJSONPB = withObject "Ping" $ \obj -> do
    service' <- obj .: "service"
    pure $ defMessage
      & P.service .~ service'

instance ToJSONPB Ping where
  toJSONPB x = object
    [ "service" .= (x^.service)
    ]
  toEncodingPB x = pairs
    [ "service" .= (x^.service)
    ]

instance FromJSON Ping where
  parseJSON = parseJSONPB

instance ToJSON Ping where
  toJSON = toAesonValue
  toEncoding = toAesonEncoding

instance FromJSONPB Pong'Extra where
  parseJSONPB = A.withObject "Pong'Extra" $ \obj -> mconcat
    [
      Pong'T <$> parseField obj "t"
    , Pong'U <$> parseField obj "u"
    ]

instance ToJSONPB Pong'Extra where
  toJSONPB (Pong'T x) = object [ "t" .= x ]
  toJSONPB (Pong'U x) = object [ "u" .= x ]
  toEncodingPB (Pong'T x) = pairs [ "t" .= x ]
  toEncodingPB (Pong'U x) = pairs [ "u" .= x ]

instance FromJSON Pong'Extra where
  parseJSON = parseJSONPB

instance ToJSON Pong'Extra where
  toJSON = toAesonValue
  toEncoding = toAesonEncoding

instance FromJSONPB Pong where
  parseJSONPB = withObject "Pong" $ \obj -> do
    status' <- obj .: "status"
    stuff' <- obj .: "stuff"
    id' <- obj .: "id"
    type'' <- obj .: "type'"
    extra' <- obj A..:? "extra"
    pure $ defMessage
      & P.status .~ status'
      & P.stuff .~ stuff'
      & P.id .~ id'
      & P.type' .~ type''
      & P.maybe'extra .~ extra'

instance ToJSONPB Pong where
  toJSONPB x = object
    [ "status" .= (x^.status)
    , "stuff" .= (x^.stuff)
    , "id" .= (x^.id)
    , "type'" .= (x^.type')
    , "extra" .= (x^.maybe'extra)
    ]
  toEncodingPB x = pairs
    [ "status" .= (x^.status)
    , "stuff" .= (x^.stuff)
    , "id" .= (x^.id)
    , "type'" .= (x^.type')
    , "extra" .= (x^.maybe'extra)
    ]

instance FromJSON Pong where
  parseJSON = parseJSONPB

instance ToJSON Pong where
  toJSON = toAesonValue
  toEncoding = toAesonEncoding

instance FromJSONPB FieldTestMessage where
  parseJSONPB = withObject "FieldTestMessage" $ \obj -> do
    testBytes' <- obj .: "testBytes"
    pure $ defMessage
      & P.testBytes .~ testBytes'

instance ToJSONPB FieldTestMessage where
  toJSONPB x = object
    [ "testBytes" .= (x^.testBytes)
    ]
  toEncodingPB x = pairs
    [ "testBytes" .= (x^.testBytes)
    ]

instance FromJSON FieldTestMessage where
  parseJSON = parseJSONPB

instance ToJSON FieldTestMessage where
  toJSON = toAesonValue
  toEncoding = toAesonEncoding

instance FromJSONPB EmptyMessage where
  parseJSONPB = withObject "EmptyMessage" $ \_ -> pure defMessage

instance ToJSONPB EmptyMessage where
  toJSONPB _ = object []
  toEncodingPB _ = pairs []

instance FromJSON EmptyMessage where
  parseJSON = parseJSONPB

instance ToJSON EmptyMessage where
  toJSON = toAesonValue
  toEncoding = toAesonEncoding

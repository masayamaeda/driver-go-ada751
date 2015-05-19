package ada751

const Ok byte = 0x00
const PacketRecieveErr byte = 0x01
const NoFinger byte = 0x02
const ImageFail byte = 0x03
const ImageMess byte = 0x06
const FeatureFail byte = 0x07
const NoMatch byte = 0x08
const NotFound byte = 0x09
const EnrollMismatch byte = 0x0a
const BadLocation byte = 0x0b
const DBRangeFail byte = 0x0c
const UploadFeatureFail byte = 0x0d
const PacketResponseFail byte = 0x0e
const UploadFail byte = 0x0f
const DeleteFail byte = 0x10
const DBClearFail byte = 0x11
const PassFail byte = 0x13
const InvalidImage byte = 0x15
const FlashErr byte = 0x18
const InvalidReg byte = 0x1a
const AddrCode byte = 0x20
const PassVerify byte = 0x21

const StartCode uint16 = 0xef01

const CommandPacket byte = 0x1
const DataPacket byte = 0x2
const AckPacket byte = 0x7
const EnddataPacket byte = 0x8

const Timeout byte = 0xff
const BadPacket byte = 0xfe

const MessageGetImage byte = 0x01
const MessageImage2tz byte = 0x02
const MessageRegModel byte = 0x05
const MessageStore byte = 0x06
const MessageLoad byte = 0x07
const MessageUpload byte = 0x08
const MessageDelete byte = 0x0c
const MessageEmpty byte = 0x0d
const MessageVerifyPassword byte = 0x13
const MessageHiSpeedSearch byte = 0x1b
const MessageTemplateCount byte = 0x1d

const DefaultTimeout uint16 = 5000

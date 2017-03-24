package main

var Colors = [...]int{
	0x5fa2bf,
	0x7bf289,
	0x8d8cdb,
	0xa3e587,
	0x281c87,
	0xe8ea56,
	0x17a55e,
	0x81e2ad,
	0x8a5ed1,
	0x80d624,
	0xbf6228,
	0xdff48b,
	0x42c401,
	0xfcc9c7,
	0xf0b2f7,
	0x2dc6b5,
	0xd64fd1,
	0x88e835,
	0xf9ac9d,
	0xea62dc,
	0x2aaf01,
	0x5dad18,
	0x93d1e2,
	0xed04c6,
	0x79b5ce,
	0x4fcdff,
	0xf75495,
	0xbe62e0,
	0x02774c,
	0x150168,
	0x425aaf,
	0xaf512f,
	0xb7d1ff,
	0xba6c0d,
	0x09d306,
	0xef32b6,
	0x72ead0,
	0xf7b2c2,
	0xc3ff7f,
	0x653af2,
	0xf74ac6,
	0x4dfc25,
	0xc0f3f9,
	0x4182bf,
	0xa71ed8,
	0xb2fca4,
	0x0e186b,
	0xa8c5ff,
	0xf7bab9,
	0x45f7a4,
	0x46d6b4,
	0xfc50c6,
	0xff89df,
	0xa2ff96,
	0xa5ffcf,
	0x5566e8,
	0xffff8c,
	0xe096dc,
	0xf492d4,
	0xf4e730,
	0xefe883,
	0x4172e2,
	0x4371ef,
	0x681ec9,
	0x6cdd87,
	0xed67f7,
	0xd64869,
	0x70e093,
	0xead131,
	0xf2ad96,
	0x360f91,
	0x095b68,
	0xd7ea91,
	0x120fdd,
	0xf97211,
	0xadffe7,
	0xaa22c9,
	0x73f48f,
	0xfcc1b3,
	0x54acce,
	0xfff1a3,
	0xffae9e,
	0xefadac,
	0x97db4e,
	0xcd8ddd,
	0x17748c,
	0xd86ecc,
	0x93ffd7,
	0x9e411f,
	0xea7033,
	0x898cd3,
	0x840204,
	0xe5375c,
	0xc41773,
	0x12a8a8,
	0xbddd63,
	0x0c00f9,
	0x3935b7,
	0xf497af,
	0x9794e0,
	0xbd1ac9,
	0xe472ea,
	0xd33da4,
	0xf9c268,
	0x649916,
	0x3890aa,
	0xef88dc,
	0x56f761,
	0x1fad55,
	0x4d84f9,
	0xf7f9a9,
	0xe5e863,
	0x74f2b5,
	0xf9fc67,
	0x933be5,
	0xedf9a7,
	0xf4c7a6,
	0xe50f04,
	0x890b00,
	0xc47a31,
	0x96aa25,
	0xd85402,
	0x23d3ca,
	0x5cf25e,
	0xff84f0,
	0xfff763,
	0xf9bbcf,
	0x02fca9,
	0xb5bff2,
	0xf9ee84,
	0x5b9bcc,
	0xdd3e00,
	0x1addc3,
	0xf2a9bc,
	0xf284f4,
	0xf9f504,
	0x8727f4,
	0xedf9a9,
	0x70b52b,
	0x61d867,
	0xf2768d,
	0xff322b,
	0xbcb5f2,
	0xeab47e,
	0xfcfca6,
	0x037551,
	0xa2e1f9,
	0xe23626,
	0x5668b7,
	0x10e5d7,
	0xf27426,
	0x0d40db,
	0xf7beed,
	0x5ee8ed,
	0xfc9cec,
	0xeaaa09,
	0xdbce3d,
	0x689de2,
	0xedc2fc,
	0xeaef56,
	0x2f58bf,
	0xfc2872,
	0x9e76ed,
	0x53d151,
	0xcf6be5,
	0xcdffb2,
	0x80f7eb,
	0xef23e2,
	0x51ef7b,
	0xf47f18,
	0x89f6f9,
	0xd1d119,
	0xe5d779,
	0xe5b86b,
	0x26c1b2,
	0x95f49d,
	0xddb771,
	0xdb6fb5,
	0x36d89f,
	0x5242aa,
	0x606ff2,
	0xce2eea,
	0xf2afab,
	0x4910ba,
	0x86e552,
	0x8e7acc,
	0x8b36db,
	0x94fca4,
	0xa0bf24,
	0x86d7e8,
	0xf2a9e1,
	0xe08d84,
	0xe228ba,
	0xf7a3bc,
	0x4ef490,
	0xe28590,
	0x896bd6,
	0xfcc2d0,
	0xba9ced,
	0x62f904,
	0xbcd81a,
	0xf298b3,
	0xd622cd,
	0x193993,
	0x902cdd,
	0x6d1cc9,
	0xa9f282,
	0x56fc4e,
	0x2cdd7f,
	0x6b42c4,
	0xffec23,
	0x5808ad,
	0x6a46ce,
	0x31e034,
	0xeda71c,
	0x230b6b,
	0x68f2d2,
	0x5c70d6,
	0x9682e5,
	0x1d72f9,
	0xdd9766,
	0x3c5ca3,
	0xf4ccff,
	0xab9df2,
	0x8a38af,
	0x8b47aa,
	0x2b24b5,
	0x0d5884,
	0x1a0f8e,
	0xdd71d4,
	0x2fa00c,
	0x7af9c4,
	0x7cf963,
	0xd863d6,
	0xd43ce8,
	0x63ffac,
	0xf9da95,
	0xf4bcb2,
	0x372c9e,
	0xc4c0f7,
	0xffdbc9,
	0xff0295,
	0xedb59e,
	0x59a4f9,
	0x1bd68b,
	0x3de5b8,
	0x649ee0,
	0xfffda3,
	0x8aeaca,
	0x169989,
	0xe9f492,
	0x380a75,
	0xed7de2,
	0x12cc2e,
	0x9ff9a1,
	0xf97522,
	0x5973f7,
	0x7f3da5,
	0xefd667,
	0xf1f446,
	0xe2818b,
	0xd912e8,
	0x6eabdd,
	0xc6097e,
	0xe589bc,
	0xf483fc,
	0xa00814,
	0x09f2b4,
	0x7d8fe0,
	0xa1c1ea,
	0xffbaed,
	0xf298c3,
	0xc1caff,
	0xe84e81,
	0x039971,
	0xf4baff,
	0xfca4c1,
	0xaaea75,
	0xc97452,
	0xeab515,
	0x02931d,
	0x34d19c,
	0x2172dd,
	0x97e0e5,
	0xa4defc,
	0x25ad20,
	0xe02cda,
	0xce5286,
	0xe59c77,
	0x040866,
	0xe5a43b,
	0x65e83a,
	0xb8fcae,
	0xeaa29f,
	0x6afce1,
	0xd3bdfc,
	0x74e8e8,
	0xf9c886,
	0xbbfc5a,
	0x4caaa9,
}

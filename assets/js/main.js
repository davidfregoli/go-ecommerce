addEventListener("pageshow", (event) => {
   const qry = new URLSearchParams(window.location.search);
   $('.price-slider').each((i, el) => {
      const vals = [qry.get('priceMin'), qry.get('priceMax') || Infinity];
      el.noUiSlider.set(vals);
   });
   $('.filter-checkbox').each((i, el) => {
      const $el = $(el);
      const name = $el.attr('name');
      const vals = qry.getAll(name);
      const checked = vals.includes($el.val());
      $el.prop('checked', checked);
   });
});

(function ($) {
   "use strict"

   // Mobile Nav toggle
   $('.menu-toggle > a').on('click', function (e) {
      e.preventDefault();
      $('#responsive-nav').toggleClass('active');
   })

   // Fix cart dropdown from closing
   $('.cart-dropdown').on('click', function (e) {
      e.stopPropagation();
   });

   /////////////////////////////////////////

   // Products Slick
   $('.products-slick').each(function () {
      var $this = $(this),
         $nav = $this.attr('data-nav');

      $this.slick({
         slidesToShow: 4,
         slidesToScroll: 1,
         autoplay: true,
         infinite: true,
         speed: 300,
         dots: false,
         arrows: true,
         appendArrows: $nav ? $nav : false,
         responsive: [{
            breakpoint: 991,
            settings: {
               slidesToShow: 2,
               slidesToScroll: 1,
            }
         },
         {
            breakpoint: 480,
            settings: {
               slidesToShow: 1,
               slidesToScroll: 1,
            }
         },
         ]
      });
   });

   // Products Widget Slick
   $('.products-widget-slick').each(function () {
      var $this = $(this),
         $nav = $this.attr('data-nav');

      $this.slick({
         infinite: true,
         autoplay: true,
         speed: 300,
         dots: false,
         arrows: true,
         appendArrows: $nav ? $nav : false,
      });
   });

   /////////////////////////////////////////

   // Product Main img Slick
   $('#product-main-img').slick({
      infinite: true,
      speed: 300,
      dots: false,
      arrows: true,
      fade: true,
      asNavFor: '#product-imgs',
   });

   // Product imgs Slick
   $('#product-imgs').slick({
      slidesToShow: 3,
      slidesToScroll: 1,
      arrows: true,
      centerMode: true,
      focusOnSelect: true,
      centerPadding: 0,
      vertical: true,
      asNavFor: '#product-main-img',
      responsive: [{
         breakpoint: 991,
         settings: {
            vertical: false,
            arrows: false,
            dots: true,
         }
      },
      ]
   });

   // Product img zoom
   var zoomMainProduct = document.getElementById('product-main-img');
   if (zoomMainProduct) {
      $('#product-main-img .product-preview').zoom();
   }

   /////////////////////////////////////////



   // Custom code
   $('.price-slider').each((i, el) => {
      const [min, max] = [+$(el).data('min'), Math.ceil(+$(el).data('max')) || 1];
      const $inputMin = $(el).siblings('.price-min').find('input');
      const $inputMax = $(el).siblings('.price-max').find('input');

      noUiSlider.create(el, {
         step: 1,
         connect: true,
         start: [$inputMin.val(), $inputMax.val()],
         range: { 'min': min, 'max': max },
      });

      el.noUiSlider.on('update', function (values, handle) {
         var value = values[handle];
         if (handle === 0) {
            $inputMin.val(value);
         } else if (handle === 1) {
            $inputMax.val(value);
         }
      });

      el.noUiSlider.on('change', () => $inputMin.change());

      $(el).parent().find('.input-number').each(function () {
         var $this = $(this),
            $input = $this.find('input'),
            $btns = $input.siblings();

         $btns.on('click', (e) => {
            const $btn = $(e.currentTarget);
            const amt = $btn.hasClass('qty-up') ? 10 : -10;
            let value = parseInt($input.val()) + amt;
            value = value < 0 ? 0 : value;
            $input.val(value);
            $input.change();
            const values = $this.hasClass('price-min') ? [value, null] : [null, value];
            el.noUiSlider.set(values);
         });
      });
   });

   $('.toggle-wishlist').on('click', async (e) => {
      const $this = $(e.currentTarget);
      const action = $this.hasClass('wishlisted') ? 'remove' : 'add';
      await $.post(`/wishlist/${action}`, { product: $this.data('product') });
      $this.toggleClass('wishlisted');
   });

   $('.store-filter').on('change', (e) => {
      if (+$('input.price-max').val() >= +$('.price-slider').data('max')) {
         $('input.price-max').prop('disabled', true);
      }
      const $form = $(e.currentTarget);
      const qs = $('.store-filter').map((i, el) => $(el).serialize()).get().join('&')
      window.location.href = `${document.location.protocol}//${document.location.host}${document.location.pathname}?${qs}`;
   });

})(jQuery);

function debounce(func, timeout = 300) {
   let timer;
   return (...args) => {
      clearTimeout(timer);
      timer = setTimeout(() => { func.apply(this, args); }, timeout);
   };
}
